# build stage
FROM golang:latest AS builder
LABEL maintainer = "hilli@github.com"
RUN mkdir -p /go/src/finance-statsd
WORKDIR /go/src/finance-statsd
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /finance-statsd .

# final stage
FROM alpine:latest
#RUN apk --no-cache add ca-certificates
COPY --from=builder /finance-statsd ./
RUN chmod +x ./finance-statsd
ENTRYPOINT ["./finance-statsd"]
# EXPOSE 3030
