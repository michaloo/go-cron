TOKEN = `cat .token`
REPO := go-cron
USER := odise
VERSION := "v0.0.1"

build:
	mkdir -p out/darwin out/linux
	GOOS=darwin go build -o out/darwin/go-cron -ldflags "-X main.build `git rev-parse --short HEAD`" bin/go-cron.go
	GOOS=linux go build -o out/linux/go-cron -ldflags "-X main.build `git rev-parse --short HEAD`" bin/go-cron.go

release: build
	rm -f out/darwin/go-cron-osx.gz
	gzip -c out/darwin/go-cron > out/darwin/go-cron-osx.gz
	rm -f out/linux/go-cron-linux.gz
	gzip -c out/linux/go-cron > out/linux/go-cron-linux.gz

