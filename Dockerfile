FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY taskify .

EXPOSE 8080

CMD ["./taskify"]
