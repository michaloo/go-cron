TOKEN = `cat .token`
REPO := go-cron
USER := odise
VERSION := "v0.0.1"

build:
	mkdir -p out/darwin out/linux
	GOOS=darwin go build -o out/darwin/go-cron -ldflags "-X main.build `git rev-parse --short HEAD`" go-cron.go
	GOOS=linux go build -o out/linux/go-cron -ldflags "-X main.build `git rev-parse --short HEAD`" go-cron.go

