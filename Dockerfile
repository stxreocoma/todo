FROM golang:1.22

WORKDIR /app

COPY web /app/web

COPY cmd ./cmd

COPY database ./database

ENV TODO_PORT=7540

ENV TODO_PASSWORD=12345

EXPOSE $TODO_PORT

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./my_app

CMD ["/app/my_app"]
