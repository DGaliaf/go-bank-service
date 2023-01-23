swag:
	swag init -g ./app/cmd/app/main.go -o ./app/docs
docker.up:
	docker-compose up -d --build
docker.down:
	docker-compose down
docker.stop:
	docker-compose stop