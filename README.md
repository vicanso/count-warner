# count-warner

count warner, when the count greater than limit, it will tigger a warn event

## API

```go
w := warner.NewWarner(time.Second, 10)
w.ResetOnWarn = true
w.On(func(key string, createdAt int64) {
  fmt.Println(key)
})
w.Inc("abcd", 1)
```
