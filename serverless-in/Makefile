.PHONY: build clean

build: clean
	env GOOS=darwin go build -ldflags="-s -w" -o bin/darwin/checkin main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/linux/checkin main.go

clean:
	rm -rf ./bin