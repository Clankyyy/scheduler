FROM golang:1.22.4 AS builder
ADD . /source
WORKDIR /source
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd/scheduler -a -o /main .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY .env ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]