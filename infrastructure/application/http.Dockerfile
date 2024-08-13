FROM golang:alpine AS builder

LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache
WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -ldflags="-s -w" -mod vendor -o /app/main cmd/http.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main /app/main
CMD ["./main"]
