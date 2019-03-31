ifeq (,$(wildcard $(STACK_CONFIG)))
    $(error STACK_CONFIG ($(STACK_CONFIG)) is not found)
endif

CODE_S3_BUCKET := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Bucket"]' -r )
CODE_S3_PREFIX := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Prefix"]' -r )
STACK_NAME := $(shell cat $(STACK_CONFIG) | jq '.["StackName"]' -r )
TEMPLATE_FILE=template.yml

all: deploy

clean:
	rm build/main

build/main: *.go
	env GOARCH=amd64 GOOS=linux go build -o build/main

sam.yml: $(TEMPLATE_FILE) build/main
	aws cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(CODE_S3_BUCKET) \
		--s3-prefix $(CODE_S3_PREFIX) \
		--output-template-file sam.yml

deploy: sam.yml
	aws cloudformation deploy \
		--template-file sam.yml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM

# 		--parameter-overrides $(PARAMETERS)