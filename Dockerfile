FROM golang:alpine as builder

WORKDIR /app

COPY go.* .
RUN go mod download
COPY . .
RUN go build -o tramo-exporter .

FROM alpine:latest
COPY --from=builder /app/tramo-exporter /bin
CMD /bin/tramo-exporter
