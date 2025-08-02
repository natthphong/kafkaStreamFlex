GOARCH=arm64 GOOS=linux go build -o goapp cmd/main.go
docker build -t test:4 .
docker run --rm -it -p 9999:9999 test:3
