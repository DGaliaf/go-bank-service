swag:
	swag init -g ./app/cmd/app/main.go -o ./app/docs
docker.up:
	docker-compose up -d --build
docker.stop:
	docker-compose stop
docker.down:
	docker-compose down
