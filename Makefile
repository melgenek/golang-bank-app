run_test:
	docker compose -f docker-compose-test.yaml up -d
	(until curl http://localhost:5432/ 2>&1 | grep '52' > /dev/null; do sleep 1; done) || true
	go clean && go test ./test/...
	docker compose down

build_docker:
	docker build -t bank_app:latest .
