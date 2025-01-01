printf "Running unit tests...\n"
printf "======================\n\n"
go test -coverprofile=coverage.out ./...

printf "\n\nChecking test coverage...\n"
printf "=========================\n\n"
go tool cover -func=coverage.out

printf "\n\nStarting the app [:8080]...\n"
printf "========================\n\n"
go run cmd/main.go