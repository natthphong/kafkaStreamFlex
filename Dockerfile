FROM alpine:latest
RUN apk --no-cache add tzdata
RUN apk --no-cache add curl
ENV TZ=Asia/Bangkok
WORKDIR /app
COPY ./goapp ./goapp
COPY pkg/sql ./sql
ENTRYPOINT [ "./goapp" ]