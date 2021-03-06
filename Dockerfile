# Compile stage
FROM golang:1.13.8 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev
# RUN go mod init github.com/manabie-com/togo
# RUN go get
RUN go build -o /server

# Final stage
FROM debian:buster

EXPOSE 5050

WORKDIR /
COPY --from=build-env /server /

CMD ["/server"]