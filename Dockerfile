FROM golang:1.18
COPY ./ /green-bot
WORKDIR /green-bot

RUN go build -o green-bot
ENTRYPOINT ["./green-bot"]
