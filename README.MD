# Query Param

[![Build Status](https://travis-ci.org/TomWright/queryparam.svg?branch=master)](https://travis-ci.org/TomWright/queryparam)
[![codecov](https://codecov.io/gh/TomWright/queryparam/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/queryparam)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/queryparam)](https://goreportcard.com/report/github.com/TomWright/queryparam)
[![Documentation](https://godoc.org/github.com/TomWright/queryparam?status.svg)](https://godoc.org/github.com/TomWright/queryparam)

Query param makes it easy to access the query parameters stored in a URL.

## Installation

```
go get -u github.com/tomwright/queryparam/v3
```

## Examples

### Parsing a simple URL
```
type GetUsersRequest struct {
	Name   string `queryparam:"name"`
	Age    string `queryparam:"age"`
	Gender string `queryparam:"gender"`
}

u, err := url.Parse("https://example.com/users?name=Tom&age=23")
if err != nil {
    panic(err)
}

req := &GetUsersRequest{}

err = queryparam.Parse(u, req)
if err != nil {
    panic(err)
}

fmt.Println(req.Name) // "Tom"
fmt.Println(req.Age) // "23"
fmt.Println(req.Gender) // ""
```

### Parsing a URL with delimited parameters
```
type GetUsersRequest struct {
	Name   []string `queryparam:"name"`
}

u, err := url.Parse("https://example.com/users?name=Tom,Jim")
if err != nil {
    panic(err)
}

req := &GetUsersRequest{}

err = queryparam.Parse(u, req)
if err != nil {
    panic(err)
}

fmt.Println(req.Name) // { "Tom", "Jim" }
```

If you want to change the delimiter you can do so by changing the value of `queryparam.Delimiter`.

The default value is `,`.

### Parsing a URL from a HTTP Request
```
var r *http.Request
// Receive r in a http handler...

type GetUsersRequest struct {
	Name   string `queryparam:"name"`
	Age    string `queryparam:"age"`
	Gender string `queryparam:"gender"`
}

req := &GetUsersRequest{}
err = queryparam.Parse(r.URL, req)
if err != nil {
    panic(err)
}

// Do something with req...
```