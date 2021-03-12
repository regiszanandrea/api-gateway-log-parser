clean:
	rm -rfv ./bin

local-test:
	go clean -testcache && go test ./...

local-integration-test:
	go clean -testcache && go test ./... -tags integration

migrate:
	docker exec -it apigatewaylog-parser /bin/sh -c "run-parts bin/migrations"

parse:
	docker exec -it apigatewaylog-parser /bin/sh -c "bin/apigateway_log_parser ${FILE_PATH}"

export-by-service:
	docker exec -it apigatewaylog-parser /bin/sh -c "bin/export_by_service ${SERVICE}"

export-by-consumer:
	docker exec -it apigatewaylog-parser /bin/sh -c "bin/export_by_consumer ${CONSUMER}"

export-metrics-by-service:
	docker exec -it apigatewaylog-parser /bin/sh -c "bin/export_metrics_by_service ${SERVICE}"

generate-coverage:
	go test -coverprofile=cover.out -coverpkg=./... ./... -tags integration;go tool cover -html=cover.out

install-hooks:
	@for f in githooks/* ; do \
  		file=`echo $$f | cut -d "/" -f2` ; \
  		echo $$file ; \
  		rm -f .git/hooks/$$file ; \
  		ln -s $(PWD)/githooks/$$file .git/hooks/$$file ; \
	done