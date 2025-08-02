FROM golang:1.24-alpine
RUN apk --no-cache add tzdata curl bash go gcc musl-dev
ENV TZ=Asia/Bangkok
WORKDIR /app
COPY ./goapp ./goapp
COPY ./config ./config
COPY pkg/sql ./sql
ENTRYPOINT [ "./goapp" ]