FROM golang:1.14 AS builder

WORKDIR /go/server
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM scratch

WORKDIR /root/
COPY --from=builder /go/server/server .
COPY --from=builder /go/server/public ./public

CMD [ "./server", "-database=redis" ]