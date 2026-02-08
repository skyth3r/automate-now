run:
	go run ./

build-linux:
	GOOS=linux GOARCH=amd64 go build ./