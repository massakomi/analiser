FROM golang:1.25.3-alpine3.22 AS builder
ENV CGO_ENABLED 0
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o hello main.go

FROM scratch
WORKDIR /build
COPY --from=builder /app/hello /build/hello
CMD ["./hello"]

