ifeq (,$(wildcard $(STACK_CONFIG)))
    $(error STACK_CONFIG ($(STACK_CONFIG)) is not found)
endif

CODE_S3_BUCKET := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Bucket"]' -r )
CODE_S3_PREFIX := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Prefix"]' -r )
STACK_NAME := $(shell cat $(STACK_CONFIG) | jq '.["StackName"]' -r )

SERVICE_DOMAIN_NAME := $(shell cat $(STACK_CONFIG) | jq '.["ServiceDomainName"]' -r )
S3_HOSTED_ZONE_ID := $(shell cat $(STACK_CONFIG) | jq '.["S3HostedZoneID"]' -r )
REGION := $(shell cat $(STACK_CONFIG) | jq '.["Region"]' -r )

TEMPLATE_FILE=template.yml

all: deploy

clean:
	rm build/main

build/main: api/*.go lambda/*.go
	env GOARCH=amd64 GOOS=linux go build -o build/main ./lambda

sam.yml: $(TEMPLATE_FILE) build/main
	aws --region $(REGION) cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(CODE_S3_BUCKET) \
		--s3-prefix $(CODE_S3_PREFIX) \
		--output-template-file sam.yml

deploy: sam.yml
	aws --region $(REGION) cloudformation deploy \
		--template-file sam.yml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
 		--parameter-overrides \
		  ServiceDomainName=$(SERVICE_DOMAIN_NAME)
	npx webpack --optimize-minimize --config ./webpack.config.js
	aws --region $(REGION) s3 sync --exact-timestamps --delete static/ s3://$(SERVICE_DOMAIN_NAME)/
	$(MAKE) sync

static/js/bundle.js: javascript/*
	npx webpack --optimize-minimize --config ./webpack.config.js

sync: static/js/bundle.js static/**
	aws --region $(REGION) s3 sync --exact-timestamps --delete static/ s3://$(SERVICE_DOMAIN_NAME)/
