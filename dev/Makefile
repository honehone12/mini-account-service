.PHONY: postgres psql test
postgres:
	docker run --rm --name NAME -e POSTGRES_PASSWORD=PW -p P:5432 postgres

psql:
	docker exec -it NAME psql -U postgres
