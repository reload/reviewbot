.PHONY: test deploy logs

NAME=ReviewBot-webhook
ENTRY_POINT=Handle
REGION=europe-west1
RUNTIME=$(shell grep '^go ' go.mod  | tr -d ' .')

export GO111MODULE=on

test: *.go
	go test ./...

deploy: test
	gcloud functions deploy $(NAME) --entry-point=$(ENTRY_POINT) --runtime=$(RUNTIME) --trigger-http --memory=128M --region=$(REGION)

logs:
	gcloud functions logs read $(NAME) --region=$(REGION) --limit=100
