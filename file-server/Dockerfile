FROM golang:1.22 as builder 
WORKDIR /build
ADD . . 
RUN go build -o app ./cmd/app

FROM ubuntu as app
WORKDIR /app
COPY --from=builder /build/app .
EXPOSE 8080
ENTRYPOINT ["./app"]
