FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./

WORKDIR /app/cmd/migrator
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/migrator

FROM scratch

COPY --from=builder /bin/migrator /migrator

ENTRYPOINT [ "/migrator" ]
