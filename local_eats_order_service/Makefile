DB_URL='postgres://postgres:1111@localhost:5433/users?sslmode=disable'
run:
	go run server/main.go
mig-up:
	migrate -database ${DB_URL} -path migrations up
mig-down:
	migrate -database ${DB_URL} -path migrations down
mig-force:
	migrate -database ${DB_URL} -path migrations force 1
mig-goto:
	migrate -database ${DB_URL} -path migrations goto 1
mig-file:
	migrate create -ext sql -dir migrations -seq payment

mod:
	go mod init composition_service
	go mod tidy
	go mod vendor
tidy:
	go mod tidy
	go mod vendor



gen-proto-common:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/transaction.proto
proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o docs




