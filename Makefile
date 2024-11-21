up:
	docker-compose up
build:
	docker-compose up --build
exec-api:
	docker-compose exec api /bin/bash
down:
	docker-compose down