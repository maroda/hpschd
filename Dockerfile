FROM alpine:latest
LABEL vendor="Sounding"
LABEL app=hpschd
LABEL org.opencontainers.image.source=https://github.com/maroda/hpschd
WORKDIR /app
COPY hpschd .
EXPOSE 9999
CMD ["./hpschd"]
