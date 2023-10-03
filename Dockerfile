FROM golang:1.21.1

RUN mkdir /app
WORKDIR /app
COPY . ./

RUN make build-all-platforms

CMD ["./builds/auth0-exporter-linux-amd64", "export", "--tls.disabled", "--log.level", "debug"]

