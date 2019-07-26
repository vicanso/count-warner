# count-warner

[![Build Status](https://img.shields.io/travis/vicanso/count-warner.svg?label=linux+build)](https://travis-ci.org/vicanso/count-warner)

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
