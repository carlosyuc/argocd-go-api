FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app/*.go ./

RUN go build -o /docker-gs-ping

EXPOSE 10101

CMD ["/docker-gs-ping"]
