BINDARY_NAME = goAuthProxy

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINDARY_NAME}-linux64 .
	GOARCH=386 GOOS=linux go build -o ./bin/${BINDARY_NAME}-linux32 .
	GOARCH=arm GOOS=linux go build -o ./bin/${BINDARY_NAME}-linux-arm .
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINDARY_NAME}-win64.exe .
	GOARCH=386 GOOS=windows go build -o ./bin/${BINDARY_NAME}-win32 .