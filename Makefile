.DEFAULT_GOAL := build

##################################
# git
##################################
GIT_URL ?= $(shell git remote get-url --push origin)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
#GIT_SHORT_COMMIT := $(shell git rev-parse --short HEAD)
TIMESTAMP := $(shell date '+%Y-%m-%d_%I:%M:%S%p')
AWS_REGION ?= us-east-2

GO_MOD_PATH ?= jimmyray.io/cdk

.PHONY: meta clean compile init check test deploy synth

meta:
	$(info    [METADATA])
	$(info    timestamp: [$(TIMESTAMP)])
	$(info    git commit: [$(GIT_COMMIT)])
	$(info    git URL: [$(GIT_URL)])
	$(info    Container image version: [$(VERSION)])
	$(info	)

compile:	clean	meta
	$(info   [COMPILE])
	go env -w GOPROXY=direct && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o cdk.bin .
	$(info	)

clean:
	-@rm cdk.bin

init:
	-@rm go.mod
	-@rm go.sum
	go mod init $(GO_MOD_PATH)
	go mod tidy

check:
	-go vet main
	-golangci-lint run

test:
	go test $(GO_MOD_PATH) -test.v

deploy:
	cdk deploy --no-rollback

synth:
	cdk synth

destroy:
	cdk destroy




