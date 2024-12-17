FROM golang:alpine AS builder

WORKDIR /build

ADD src/go.mod .

COPY . .

RUN go build -o aStar src/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/aStar /build/aStar

CMD [". /aStar"]