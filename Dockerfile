FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o main main.go

# Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY wait-for.sh .
COPY start.sh .
COPY db/migration ./db/migration
RUN chmod +x wait-for.sh
RUN chmod +x start.sh

EXPOSE 1234
CMD [ "/app/main" ]