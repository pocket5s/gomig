.PHONY: all test clean run cover 
GOCMD=go

test:
	@if [ ! -f /go/bin/gotest ]; then \
		echo "installing gotest..."; \
		go get github.com/rakyll/gotest; \
		go install github.com/rakyll/gotest; \
	fi 
	gotest -v .

cover:
	gotest -covermode=count -coverpkg=. -coverprofile=profile.cov ./... fmt
	go tool cover -func=profile.cov
	go tool cover -html=profile.cov -o=coverage.html

clean:
	go clean
	rm gomig

run:
	go run cmd/main.go

build:
	@-$(MAKE) -s clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o gomig main.go
	chmod +x gomig
