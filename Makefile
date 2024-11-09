migrate-up:
	goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy" up
migrate-down:
	goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy" down
reset-db: migrate-down migrate-up

psql:
	psql "postgres://postgres:postgres@localhost:5432/chirpy"

asdf: