# count-warner

[![Build Status](https://github.com/vicanso/count-warner/workflows/Test/badge.svg)](https://github.com/vicanso/count-warner/actions)

Count warner, when the count greater than limit, it will emit a warn event

## API

```go
w := warner.NewWarner(time.Second, 10)
w.ResetOnWarn = true
w.On(func(key string, count int) {
  fmt.Println(key)
})
w.Inc("abcd", 1)
```
