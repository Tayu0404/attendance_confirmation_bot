FROM golang:1.12.5-stretch

MAINTAINER Tayu

WORKDIR /go

# ADD . /go

RUN go get -u github.com/bwmarrin/discordgo \
              github.com/joho/godotenv

COPY ./src /go/

RUN go build main.go
