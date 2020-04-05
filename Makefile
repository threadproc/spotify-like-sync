.PHONY: main

build: main.go
	GOOS=linux GOARCH=amd64 go build -o likesync_linux_amd64 main.go
	build-lambda-zip -output lambda.zip likesync_linux_amd64

upload: build
	aws lambda update-function-code --function-name like-sync --zip-file fileb://lambda.zip

test:
	go run main.go test