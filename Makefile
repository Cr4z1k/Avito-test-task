build:
	sudo docker-compose build
run-w-tests:
	sudo RUN_TESTS=true docker-compose up
run-w-tests-bg:
	sudo RUN_TESTS=true docker-compose up -d
run:
	sudo docker-compose up
run-bg:
	sudo docker-compose up -d
stop:
	sudo docker-compose stop
down:
	sudo docker-compose down
