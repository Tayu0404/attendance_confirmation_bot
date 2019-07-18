FROM golang:stretch as builder

RUN go get github.com/bwmarrin/discordgo \
		   github.com/jmoiron/sqlx \
		   github.com/carlescere/scheduler

RUN go get -u github.com/go-sql-driver/mysql

WORKDIR /go/src/github.com/Tayu0404/attendance_rec/discord_bot/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine

ENV DOCKERIZE_VERSION v0.6.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN apk --no-cache add --virtual=wget wget \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& apk del wget	
COPY --from=builder /go/src/github.com/Tayu0404/attendance_rec/discord_bot/app .
COPY ./schedule.json .

ENTRYPOINT ["dockerize", "-wait", "tcp://db:3306"]
CMD ["./app"]
