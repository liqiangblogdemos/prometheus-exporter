FROM golang:1.14 AS build-env

ADD . /go/src/github.com/liqiangblogdemos/prometheus-exporter
WORKDIR /go/src/github.com/liqiangblogdemos/prometheus-exporter

RUN GOOS=linux GO111MODULE=off go build -o prom-exporter main.go

FROM busybox:glibc
COPY --from=build-env /go/src/github.com/liqiangblogdemos/prometheus-exporter/prom-exporter /prom-exporter
RUN chmod +x /prom-exporter

# PORT
EXPOSE 8080

CMD ["/prom-exporter"]
