FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 go build -o /api ./cmd/api
RUN CGO_ENABLED=0 go build -o /cli ./cmd/cli

FROM alpine:3.20

COPY --from=builder /api /api
COPY --from=builder /cli /cli
EXPOSE 8080
CMD ["/api"]
