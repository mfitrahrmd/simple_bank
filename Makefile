db_username = dev
db_password = dev
db_host = localhost
db_port = 5432
db_name = simple_bank

dsn = postgres://${db_username}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=disable

loc = ./bin
initialize:
	mkdir -p ${loc}
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz --directory ${loc}

db-start:
	docker run -dt --name pg -p ${db_port}:${db_port} -e POSTGRES_USER=${db_username} -e POSTGRES_PASSWORD=${db_password} -e TZ=Asia/Jakarta -e PGTZ=Asia/Jakarta postgres:alpine

db-create:
	docker exec pg createdb --username=${db_username} ${db_name}

db-drop:
	docker exec pg dropdb --username=${db_username} ${db_name}

migrate-up:
	${loc}/migrate -source file:./database/migrations -database ${dsn} -verbose up

migrate-down:
	${loc}/migrate -source file:./database/migrations -database ${dsn} -verbose down

.PHONY: initialize startdb createdb