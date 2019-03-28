include ./helm/secrets/values.mk

NAME:=k8s-rds
EXAMPLE_NAMESPACE:=c12345-test
EXAMPLE_NAME=c12345-name

run:
	go run main.go

debug-publish: build push get delete

build:
	docker build -t gcr.io/totvscloud104/k8s-rds:latest .

push:
	docker push gcr.io/totvscloud104/k8s-rds:latest

get:
	k -n ${NAME} get pod -l name=${NAME}

delete:
	k -n ${NAME} delete pod -l name=${NAME}

logs:
	k -n ${NAME} logs $(shell k -n ${NAME} get pod -l name=${NAME} -o jsonpath='{.items[0].metadata.name}')

# Example
example-oracle:
	kubectl -n ${EXAMPLE_NAMESPACE} apply -f ./examples/oracle.yaml

# example-postgres:
# 	kubectl -n ${EXAMPLE_NAMESPACE} apply -f ./examples/postgres.yaml

example-get:
	k -n ${EXAMPLE_NAMESPACE} get database ${EXAMPLE_NAME}

example-delete:
	kubectl -n ${EXAMPLE_NAMESPACE} delete database ${EXAMPLE_NAME}

.PHONY: run example-oracle example-oracle-delete example-postgres example-postgres-delete logs
