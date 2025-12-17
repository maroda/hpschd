FROM golang:1.25.4-alpine
LABEL vendor="Sounding"
EXPOSE 9999
WORKDIR /go/src/hpschd/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/hpschd
ENTRYPOINT ["/bin/hpschd"]
