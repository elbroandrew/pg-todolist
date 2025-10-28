FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./

WORKDIR /app/cmd/task_service
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/task_service

FROM alpine:latest

# add curl for healthchecks
RUN apk add --no-cache curl

COPY --from=builder /bin/task_service /task_service

EXPOSE 8081

ENTRYPOINT [ "/task_service" ]
