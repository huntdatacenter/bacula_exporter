FROM golang:1.26 AS builder
WORKDIR /go/src/github.com/huntdatacenter/bacula_exporter
ENV GOPATH=/go/
COPY . .
RUN make clean && make bacula_exporter

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/huntdatacenter/bacula_exporter/bacula_exporter /app/bacula_exporter
EXPOSE 33407
CMD ["/app/bacula_exporter", "--config", "/config/bacula_exporter.knf"]
