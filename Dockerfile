FROM golang:1.15 AS builder

WORKDIR /go/server
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM drone/ca-certs

WORKDIR /root/
COPY --from=builder /go/server/server .
COPY --from=builder /go/server/public ./public

EXPOSE 6443

CMD [ "./server", "-database=redis" ]