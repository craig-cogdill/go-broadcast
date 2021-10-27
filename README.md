# go-broadcast

A very basic broadcast / subscription library for handling one-to-many queuing in Golang.

### Example usage
```
// main.go
b := broadcast.New()

go func() {
   subscription := b.Subscribe()
   msg := <-subscription.Queue()
}()

b.Broadcast("Hello World")

```

### Running tests

```
make test
```

### Code coverage

```
make coverage
```
