FROM golang:alpine
WORKDIR /build
COPY src/main.go .
RUN go build -o hello main.go
CMD [". /hello"]