.PHONY: docker-seed
docker-seed:
	./scripts/seed-data.sh "docker"

.PHONY: seed
seed:
	./scripts/seed-data.sh "local"

.PHONY: test-seed
test-seed:
	./scripts/seed-data.sh "test"