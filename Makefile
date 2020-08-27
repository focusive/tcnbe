newdb:
	docker exec -i some-mysql sh -c 'exec mysql -uroot -p"my-secret-pw"' < "./initial.sql"
migrate:
	go run main.go --migrate