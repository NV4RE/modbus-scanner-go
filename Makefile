APP_NAME=modbus-scanner

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)-$(shell date +%Y-%m-%d-%H-%M).exe .