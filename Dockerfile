FROM golang:1.19-bullseye 

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /fampay-video-library-exe

EXPOSE 8080

CMD [ "/fampay-video-library-exe" ]