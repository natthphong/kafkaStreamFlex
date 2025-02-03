FROM alpine:latest
RUN apk --no-cache add tzdata curl bash go gcc musl-dev
ENV TZ=Asia/Bangkok
WORKDIR /app
COPY ./goapp ./goapp
COPY pkg/sql ./sql
ENTRYPOINT [ "./goapp" ]