FROM golang:latest as builder

RUN mkdir app
WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN go build -o bin/server ./cmd

FROM busybox

RUN mkdir app

COPY --from=builder /app/bin/* /app

CMD ["./app/server"]
