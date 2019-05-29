FROM golang:stretch as builder

RUN go get -u github.com/bwmarrin/discordgo \
              github.com/joho/godotenv

WORKDIR /go/src/discord_bot/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine

ENV DOCKERIZE_VERSION v0.6.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN apk update && apk add -y wget \
	&& wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& apk del --purge wget
	
COPY --from=builder /go/src/discord_bot/app .
COPY .env .

ENTRYPOINT ["dockerize", "-wait", tcp://db:3306]
RUN ls -a
CMD ["./app"]
