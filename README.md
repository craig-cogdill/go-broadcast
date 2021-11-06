# go-broadcast

### Example usage
```
// main.go
b := broadcast.New()

go func() {
   subscription := b.Subscribe()
   msg := <-subscription.Queue()
   log.Print(msg.(string)) // "Hello World"
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

Interactive hotspot coverage report:
```
make coverage-html
```
