
all: 
	GOPATH=$(shell pwd) go build src/main.go

run: all 
	sudo systemctl stop fan_ctl
	sudo ./main
	sudo systemctl start fan_ctl