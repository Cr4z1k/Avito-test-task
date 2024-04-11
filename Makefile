build:
	docker-compose build
run-w-tests:
	RUN_TESTS=true docker-compose up
run:
	docker-compose up
stop:
	docker-compose stop
down:
	docker-compose down
