FROM private-repo.tar-server.com/alpine-go:1.24.0
RUN apk --no-cache add tzdata curl bash go gcc musl-dev
ENV TZ=Asia/Bangkok
WORKDIR /app
COPY ./goapp ./goapp
COPY go.mod go.sum ./
RUN go mod download

COPY ./config ./config
COPY pkg/sql ./sql
ENTRYPOINT [ "./goapp" ]