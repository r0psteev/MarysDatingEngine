FROM golang:alpine as builder
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN go build .

FROM alpine
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]
