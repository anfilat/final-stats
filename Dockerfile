FROM golang:1.15-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/symo /app/cmd/symo/


FROM alpine:latest

RUN apk add --no-cache coreutils && apk add --no-cache sysstat

WORKDIR /root/
COPY --from=builder /app/symo .
EXPOSE 8000

CMD ["./symo"]
