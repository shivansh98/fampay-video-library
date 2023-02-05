# fampay-video-library

## Install golang
Here used golang version 1.19.3

## Install Docker
for building and deploying docker image

## Run
To run the application without docker , can simply build the application with `go build` command
and run the .exe file. The server would start be running at localhost:8080 port.

## Run in Docker Container
build the docker image using below command
`docker build -t  fampay-video-library .`

after the build got completed run the docker run command
`docker run -p 8080:8080 fampay-video-library`

## Request structure of REST Apis
We have 2 APIs 
1. /get-videos called with GET method 
`curl --location --request GET 'http://localhost:8080/get-videos'`

2. /search called with POST method
``curl --location --request POST 'http://localhost:8080/search' \
--header 'Content-Type: application/json' \
--data-raw '{
    "text":"Live cricket in india"
}'``

## Cron
We have a go routine that runs recursively with delay of 1 min to fetch data from youtube data api and insert it in MongoDB atlas
