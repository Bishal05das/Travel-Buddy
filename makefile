migration_up:
	 migrate -path ./migrations -database "postgresql://postgres:password@localhost:5435/mydb?sslmode=disable" -verbose up
