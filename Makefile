EXECUTABLE := server

default:
	go build -o $(EXECUTABLE) main.go

test:
	go test -v ./...

coverage:
	go test -v -cover ./...

clean:
	rm $(EXECUTABLE)