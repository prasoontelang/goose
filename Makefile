.PHONY: dist
dist:
	@mkdir -p ./bin
	@rm -f ./bin/*
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/goose-darwin64       ./cmd/goose
	GOOS=linux   GOARCH=amd64 go build -o ./bin/goose-linux64        ./cmd/goose
	GOOS=linux   GOARCH=386   go build -o ./bin/goose-linux386       ./cmd/goose
	GOOS=windows GOARCH=amd64 go build -o ./bin/goose-windows64.exe  ./cmd/goose
	GOOS=windows GOARCH=386   go build -o ./bin/goose-windows386.exe ./cmd/goose

	GOOS=darwin  GOARCH=amd64 go build -o ./bin/geese-darwin64       ./cmd/geese
	GOOS=linux   GOARCH=amd64 go build -o ./bin/geese-linux64        ./cmd/geese
	GOOS=linux   GOARCH=386   go build -o ./bin/goose-linux386       ./cmd/geese
	GOOS=windows GOARCH=amd64 go build -o ./bin/geese-windows64.exe  ./cmd/geese
	GOOS=windows GOARCH=386   go build -o ./bin/geese-windows386.exe ./cmd/geese
.PHONY: vendor
vendor:
	mv _go.mod go.mod
	mv _go.sum go.sum
	GO111MODULE=on go build -o ./bin/goose ./cmd/goose
	GO111MODULE=on go mod vendor && GO111MODULE=on go mod tidy
	mv go.mod _go.mod
	mv go.sum _go.sum
