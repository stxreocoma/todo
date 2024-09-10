FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV TODO_PORT=7540

ENV TODO_PASSWORD=12345

EXPOSE $TODO_PORT

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app ./cmd

CMD ["/my_app"]