BIN=orders
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
ifeq ($(OS), Windows_NT)
	BIN=orders.exe
endif

all:
	@go clean
	@echo "Building Orders Server"
	@echo "Compiling"
	@cd ./cmd/;\
	GOOS=${GOOS} GOARCH=${GOARCH} go version
	@cd ./cmd/;\
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../bin/${BIN}
	@echo "Done"