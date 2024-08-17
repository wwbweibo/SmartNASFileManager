FROM golang:1.22 as builder 
WORKDIR /build
ADD . . 
ENV HTTPS_PROXY=http://192.168.163.65:10080
RUN go build -o app ./cmd/app

FROM ubuntu as app
WORKDIR /app
COPY --from=builder /build/app .
EXPOSE 8080
ENTRYPOINT ["./app"]
