# build
FROM golang:alpine as BUILD
WORKDIR /src
ADD ./ /src
ADD ./internal/config /app/config
RUN go build -o /app


# publish
FROM alpine

ENV TARGET=""
ENV SERVICE=""
ENV APP_VERSION=""
ENV ServerID=""
ENV TZ=Asia/Taipei

WORKDIR /app
RUN apk update && apk add tzdata ffmpeg

COPY --from=BUILD ["/app","/app"]
COPY --from=BUILD ["/app/config", "/app/config"]
COPY --from=BUILD ["/src/migration", "/app/migration"]

ENTRYPOINT /app/botgpt -e ${TARGET} -s ${SERVICE} -v ${APP_VERSION} -i ${ServerID}
