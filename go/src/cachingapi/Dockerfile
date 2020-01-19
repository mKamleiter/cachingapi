FROM golang
RUN go get github.com/bmizerany/pat && go get github.com/mattn/go-sqlite3
RUN cd / && \
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=DE/ST=Testing/L=Testing/O=Testing/OU=Testing/CN=localhost"
COPY main.go /go/src/cachingapi/main.go
RUN cd /go/src/cachingapi/ && \
    go build && go install && \
    mv /go/bin/cachingapi /cachingapi

CMD ["/cachingapi"]
