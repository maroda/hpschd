FROM golang:1.15.1-alpine3.12
LABEL version="1.2.0"
LABEL vendor="Sounding"
EXPOSE 9999
WORKDIR /go/src/hpschd/
COPY . .
RUN apk add --no-cache git
RUN go get github.com/gorilla/mux
RUN go get github.com/rs/zerolog
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go build -o /bin/hpschd
ENTRYPOINT ["/bin/hpschd"]
