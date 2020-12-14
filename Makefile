.PHONY : all
all : build

build:
	GOOS=linux GOARCH=amd64 go build -o bin/ldaptokenauth cmd/ldaptokenauth/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/ldaptokenauth-server cmd/ldaptokenauth-server/main.go
