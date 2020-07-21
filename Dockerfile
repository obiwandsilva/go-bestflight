FROM golang:1.14.3-alpine3.11 AS builder
WORKDIR /tmp/go-bestflight
COPY . .
RUN ["go", "mod", "download"]
RUN ["go", "build", "-o", "bestflight","cmd/bestflight/main.go"]

FROM alpine:3.11
WORKDIR /opt/bestflight
COPY --from=builder /tmp/go-bestflight/bestflight .
COPY --from=builder /tmp/go-bestflight/input.csv .
EXPOSE 5000
CMD [ "./bestflight", "input.csv", "5000" ]