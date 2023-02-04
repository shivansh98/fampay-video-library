FROM golang:1.19-bullseye 

WORKDIR /fampay-video-library

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o /fampay-video-library

EXPOSE 8080

CMD [ "/fampay-video-library" ]