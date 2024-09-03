FROM golang:1.22

WORKDIR /app

COPY auth ./auth

COPY cmd ./

COPY database ./database

COPY handlers ./handlers

COPY models ./models

COPY validation ./validation

COPY go.mod go.sum ./

RUN go mod download

RUN go get github.com/stxreocoma/todo/utils

ENV TODO_PORT=7540

ENV TODO_PASSWORD=12345

EXPOSE $TODO_PORT

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./my_app

CMD ["./my_app"]
