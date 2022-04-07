FROM golang:1.16

WORKDIR /app

COPY . . 

RUN go build .

COPY go.mod .
COPY go.sum .


CMD ["./pingoceantask"]