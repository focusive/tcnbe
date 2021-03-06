db:
	docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql
newdb:
	docker exec -i some-mysql sh -c 'exec mysql -uroot -p"my-secret-pw"' < "./initial.sql"
migrate:
	go run main.go --migrate
run:
	DB_CONN="ktb@ktbserver:Passw0rd@tcp(ktbserver.mysql.database.azure.com)/thaichana?charset=utf8&parseTime=True&loc=Local" go run main.go
test:
	go test -v ./...
acceptance:
	restcli tests/suite.http