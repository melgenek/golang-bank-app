FROM golang:1.17.0 as builder
COPY . /opt/app
WORKDIR /opt/app
RUN go mod download
RUN CGO_ENABLED=0 go build -o bank_app ./src

FROM alpine:3.14.3
COPY --from=builder /opt/app/bank_app /opt/app/bank_app
COPY --from=builder /opt/app/config.yaml /opt/app/config.yaml
WORKDIR /opt/app
EXPOSE 8000
ENTRYPOINT ["./bank_app"]
