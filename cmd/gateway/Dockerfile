FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./

WORKDIR /app/cmd/gateway
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/gateway

FROM scratch

COPY --from=builder /bin/gateway /gateway

EXPOSE 8080

ENTRYPOINT [ "/gateway" ]
