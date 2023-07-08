FROM golang:1.20-alpine3.17 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app/demandSvc

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/demandSvc /app/
RUN chmod +x /app/demandSvc
CMD [ "/app/demandSvc" ]
