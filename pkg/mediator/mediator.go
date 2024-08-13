package mediator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
)

// RequestHandlerFunc is a continuation for the next task to execute in the pipeline
type RequestHandlerFunc func(ctx context.Context) (interface{}, error)

// PipelineBehavior is a Pipeline behavior for wrapping the inner handler.
type PipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (interface{}, error)
}

type RequestHandler[TRequest any, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type RequestHandlerFactory[TRequest any, TResponse any] func() RequestHandler[TRequest, TResponse]

type NotificationHandler[TNotification any] interface {
	Handle(ctx context.Context, notification TNotification) error
}

type NotificationHandlerFactory[TNotification any] func() NotificationHandler[TNotification]

var requestHandlersRegistrations = map[reflect.Type]interface{}{}
var notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
var pipelineBehaviours []interface{}

type Unit struct{}

func registerRequestHandler[TRequest any, TResponse any](handler any) error {
	var request TRequest
	requestType := reflect.TypeOf(request)

	_, exist := requestHandlersRegistrations[requestType]
	if exist {
		// each request in request/response strategy should have just one handler
		return fmt.Errorf("registered handler already exists in the registry for message %s", requestType.String())
	}

	requestHandlersRegistrations[requestType] = handler

	return nil
}

func RegisterRequestHandler[TRequest any, TResponse any](handler RequestHandler[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](handler)
}

func RegisterRequestHandlerFactory[TRequest any, TResponse any](factory RequestHandlerFactory[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](factory)
}

func RegisterRequestPipelineBehaviors(behaviours ...PipelineBehavior) error {
	for _, behavior := range behaviours {
		behaviorType := reflect.TypeOf(behavior)

		existsPipe := existsPipeType(behaviorType)
		if existsPipe {
			return fmt.Errorf("registered behavior already exists in the registry")
		}

		pipelineBehaviours = append(pipelineBehaviours, behavior)
	}

	return nil
}

func registerNotificationHandler[TEvent any](handler any) error {
	var event TEvent
	eventType := reflect.TypeOf(event)

	handlers, exist := notificationHandlersRegistrations[eventType]

	if !exist {
		notificationHandlersRegistrations[eventType] = []interface{}{handler}
		return nil
	}

	notificationHandlersRegistrations[eventType] = append(handlers, handler)
	return nil
}

func RegisterNotificationHandler[TEvent any](handler NotificationHandler[TEvent]) error {
	return registerNotificationHandler[TEvent](handler)
}

func RegisterNotificationHandlerFactory[TEvent any](factory NotificationHandlerFactory[TEvent]) error {
	return registerNotificationHandler[TEvent](factory)
}

func RegisterNotificationHandlers[TEvent any](handlers ...NotificationHandler[TEvent]) error {
	if len(handlers) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range handlers {
		err := RegisterNotificationHandler(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterNotificationHandlersFactories[TEvent any](factories ...NotificationHandlerFactory[TEvent]) error {
	if len(factories) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range factories {
		err := RegisterNotificationHandlerFactory[TEvent](handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearRequestRegistrations() {
	requestHandlersRegistrations = map[reflect.Type]interface{}{}
}

func ClearNotificationRegistrations() {
	notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
}

func buildRequestHandler[TRequest any, TResponse any](handler any) (RequestHandler[TRequest, TResponse], bool) {
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		factory, ok := handler.(RequestHandlerFactory[TRequest, TResponse])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

// Send the request to its corresponding request handler.
func Send[TRequest any, TResponse any](ctx context.Context, request TRequest) (TResponse, error) {
	requestType := reflect.TypeOf(request)
	var response TResponse
	handler, ok := requestHandlersRegistrations[requestType]
	if !ok {
		return *new(TResponse), fmt.Errorf("no handler for request %T", request)
	}

	handlerValue, ok := buildRequestHandler[TRequest, TResponse](handler)
	if !ok {
		return *new(TResponse), fmt.Errorf("handler for request %T is not a Handler", request)
	}

	if len(pipelineBehaviours) > 0 {
		finalHandler := func(ctx context.Context) (interface{}, error) {
			return handlerValue.Handle(ctx, request)
		}

		for i := len(pipelineBehaviours) - 1; i >= 0; i-- {
			pipe := pipelineBehaviours[i]
			next := finalHandler

			finalHandler = func(ctx context.Context) (interface{}, error) {
				return pipe.(PipelineBehavior).Handle(ctx, request, next)
			}
		}

		result, err := finalHandler(ctx)
		if err != nil {
			return *new(TResponse), fmt.Errorf("error handling request: %w", err)
		}

		return result.(TResponse), nil
	} else {
		res, err := handlerValue.Handle(ctx, request)
		if err != nil {
			return *new(TResponse), fmt.Errorf("error handling request: %w", err)
		}

		response = res
	}

	return response, nil
}

func validateNotificationHandler(handler interface{}) error {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Ptr && handlerType.Kind() != reflect.Struct {
		return fmt.Errorf("invalid handler: must be a struct or pointer to struct, got %T", handler)
	}

	method, ok := handlerType.MethodByName("Handle")
	if !ok {
		return fmt.Errorf("invalid handler: must have a Handle method")
	}

	if method.Type.NumIn() != 3 {
		return fmt.Errorf("invalid Handle method: must have 2 input parameters")
	}

	if method.Type.In(1) != contextType {
		return fmt.Errorf("invalid Handle method: first parameter must be context.Context")
	}

	if method.Type.In(2).Kind() == reflect.Ptr {
		return fmt.Errorf("invalid Handle method: second parameter cannot be a pointer")
	}

	if method.Type.In(2).Kind() == reflect.Interface {
		return fmt.Errorf("invalid Handle method: second parameter cannot be an interface")
	}

	if method.Type.NumOut() != 1 {
		return fmt.Errorf("invalid Handle method: must have 1 output parameter")
	}

	if method.Type.Out(0) != errorType {
		return fmt.Errorf("invalid Handle method: output must be an error")
	}

	return nil
}

func Publish[TNotification any](ctx context.Context, notifications ...TNotification) error {
	for _, notification := range notifications {
		eventType := reflect.TypeOf(notification)
		handlers, ok := notificationHandlersRegistrations[eventType]

		if !ok {
			return fmt.Errorf("no handlers found for event type %s", eventType)
		}

		for _, handler := range handlers {
			if err := validateNotificationHandler(handler); err != nil {
				return err
			}

			handlerValue := reflect.ValueOf(handler)
			method := handlerValue.MethodByName("Handle")

			if !method.IsValid() {
				return fmt.Errorf("handler does not have a Handle method")
			}

			// Call the Handle method using reflection
			results := method.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(notification),
			})

			// Check for error returned by Handle
			if errInterface := results[0].Interface(); errInterface != nil {
				if err, ok := errInterface.(error); ok {
					return err
				}
				return fmt.Errorf("unexpected error type returned from handler")
			}
		}
	}

	return nil
}

func PublishAsync[TNotification any](ctx context.Context, notifications ...TNotification) error {
	for _, notification := range notifications {
		eventType := reflect.TypeOf(notification)
		handlers, ok := notificationHandlersRegistrations[eventType]

		if !ok {
			return fmt.Errorf("no handlers found for event type %s", eventType)
		}

		for _, handler := range handlers {
			if err := validateNotificationHandler(handler); err != nil {
				return err
			}

			go func(handler interface{}) {
				handlerValue := reflect.ValueOf(handler)
				method := handlerValue.MethodByName("Handle")

				if !method.IsValid() {
					return
				}

				// Call the Handle method using reflection
				results := method.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					reflect.ValueOf(notification),
				})

				// Check for error returned by Handle
				if errInterface := results[0].Interface(); errInterface != nil {
					if err, ok := errInterface.(error); ok {
						fmt.Printf("error: %s", err.Error())
					}
				}
			}(handler)
		}
	}
	return nil
}

func existsPipeType(p reflect.Type) bool {
	for _, pipe := range pipelineBehaviours {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}
