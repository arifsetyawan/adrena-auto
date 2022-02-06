.PHONY: build clean

build: clean
	env GOOS=darwin go build -ldflags="-s -w" -o bin/darwin/checkin In/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/linux/checkin In/main.go
	env GOOS=windows go build -ldflags="-s -w" -o bin/windows/checkin In/main.go

clean:
	rm -rf ./bin