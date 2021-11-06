EXECUTABLE := server
COVERAGE_REPORT := coverage.out

default:
	go build -o $(EXECUTABLE) main.go

test:
	go test -v ./...

coverage:
	go test -v -cover ./...

coverage-html:
	go test -v -coverprofile=$(COVERAGE_REPORT) ./...
	go tool cover -html=$(COVERAGE_REPORT)

clean:
	rm $(EXECUTABLE) $(COVERAGE_REPORT)