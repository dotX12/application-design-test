# Application design

### Logs format
```
2024-08-13T14:05:59+03:00 TRC internal\presentation\http\v1\handler\create_order.go:26 > Received request to CreateOrderHandler spanId=c9ff6fa7255542bb traceId=66bb3e17076a3fd1f1552e18bd70264b
2024-08-13T14:05:59+03:00 TRC internal\application\pipeline\logging.go:29 > Received request for processing in mediator request={"email":"foo@bar.com","from":"2025-01-02T00:00:00Z","hotel_id":"string","room_id":"string","to":"2025-01-02T00:00:00Z"} spanId=5f6108f795686b7e traceId=66bb3e17076a3fd1f1552e18bd70264b type=AddOrderCommand
2024-08-13T14:05:59+03:00 TRC internal\application\command\create_order.go:43 > AddOrderCommandHandler.Handle spanId=5f6108f795686b7e traceId=66bb3e17076a3fd1f1552e18bd70264b
2024-08-13T14:05:59+03:00 TRC internal\application\pipeline\logging.go:52 > Request processed successfully in mediator request={"email":"foo@bar.com","from":"2025-01-02T00:00:00Z","hotel_id":"string","room_id":"string","to":"2025-01-02T00:00:00Z"} response={"From":"2025-01-02T00:00:00Z","HotelID":{},"ID":{},"RoomID":{},"To":"2025-01-02T00:00:00Z","UserEmail":{}} spanId=5f6108f795686b7e traceId=66bb3e17076a3fd1f1552e18bd70264b type=AddOrderCommand
2024-08-13T14:05:59+03:00 TRC internal\application\event\order_created.go:30 > We send a email about the hotel reservation... created_at="2024-08-13 14:05:59.8012522 +0300 MSK m=+200.421637801" email=foo@bar.com event_id=f069d186-cdd0-4613-983d-4d7dc33006bd from=2025-01-02T00:00:00Z order_id=fbf3de2a-8cbe-48c6-ae04-faf46f84802d spanId=5f6108f795686b7e to=2025-01-02T00:00:00Z topic=booking.OrderCreated traceId=66bb3e17076a3fd1f1552e18bd70264b
2024-08-13T14:05:59+03:00 INF pkg\logging\fiberzerolog\fiber.go:80 > Success bytesReceived=149 bytesSent=144 ip=127.0.0.1 latency="89.7µs" method=POST spanId=c9ff6fa7255542bb status=201 traceId=66bb3e17076a3fd1f1552e18bd70264b ua="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36" url=/v1/orders
```

## Docs

![image](https://github.com/user-attachments/assets/7277ce2b-d540-47d5-bcd8-db71b9dcb9bc)
![image](https://github.com/user-attachments/assets/eaf64939-fac6-46af-a94d-ca26784478aa)

## Metrics
![image](https://github.com/user-attachments/assets/debbd7b0-c607-4547-b151-08fa5e95612f)

```sh
swag init -g cmd/http.go --instanceName v1 --parseDependency --collectionFormat multi
```
