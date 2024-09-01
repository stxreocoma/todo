FROM ubuntu:latest

WORKDIR /app

COPY main ./

COPY web /app/web

ENV TODO_PORT=7540

ENV TODO_PASSWORD=12345

EXPOSE $TODO_PORT

CMD ["/app/main"]
