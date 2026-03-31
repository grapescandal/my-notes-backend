FROM golang:1.25-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Then copy source
COPY . .

RUN go build -o /my-note

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /my-note /my-note
EXPOSE 8080
ENTRYPOINT ["/my-note"]
