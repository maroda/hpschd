FROM alpine:latest
LABEL vendor="Sounding"
LABEL app=hpschd
LABEL org.opencontainers.image.source=https://github.com/maroda/hpschd
WORKDIR /app
COPY hpschd .
COPY public/ ./public/
EXPOSE 9999
CMD ["./hpschd"]
