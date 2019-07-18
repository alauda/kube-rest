
all: test

# Run tests
test: fmt vet
	go test ./pkg/... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...
