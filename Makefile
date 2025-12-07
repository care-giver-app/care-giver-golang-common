test:
	go test -short -coverprofile cover.out ./...

test-report: test
	go tool cover -html=cover.out

lint: 
	golangci-lint run