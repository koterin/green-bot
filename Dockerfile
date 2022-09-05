FROM golang:1.18 AS builder
WORKDIR /build

COPY . /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/main.go

FROM alpine:3.16 as server
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main /green-bot/

WORKDIR /green-bot

ENTRYPOINT [ "./main" ]
