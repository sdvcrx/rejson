# rejson

> Reshape JSON model. Based on [gjson](https://github.com/tidwall/gjson) .

[![Github Action](https://github.com/sdvcrx/rejson/workflows/Go/badge.svg)](https://github.com/sdvcrx/rejson/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/sdvcrx/rejson)](https://goreportcard.com/report/github.com/sdvcrx/rejson)
[![codecov](https://codecov.io/gh/sdvcrx/rejson/branch/master/graph/badge.svg?token=WJVJ0WRX3C)](https://codecov.io/gh/sdvcrx/rejson)
[![Go.dev](https://pkg.go.dev/badge/github.com/sdvcrx/rejson)](https://pkg.go.dev/github.com/sdvcrx/rejson)

## Install

Run `go get`:

```sh
go get -u github.com/sdvcrx/rejson
```

## Usage

```go
package main

import (
	"log"

	"github.com/sdvcrx/rejson"
)

type User struct {
	Name string `rejson:"data.name"`
}

func main() {
	jsonString := `{
		"code":0,
		"msg": null,
		"data":{"name":"John"}
	}`

	u := User{}
	rejson.Unmarshal(jsonString, &u)
	log.Printf("%#v", u)
	// u => User{Name:John}
}
```

## Performance

Benchmark cases: [bench_test.go](https://github.com/sdvcrx/rejson/blob/master/bench_test.go)

```
BenchmarkUnmarshalReJSON-4                818850              1443 ns/op              72 B/op            4 allocs/op
BenchmarkUnmarshalGJSONGet-4             1512412               798 ns/op               0 B/op            0 allocs/op
BenchmarkUnmarshalEncodingJSON-4          312958              3785 ns/op             464 B/op           11 allocs/op
```

## License

MIT
