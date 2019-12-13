# Query Param

[![Build Status](https://travis-ci.org/TomWright/queryparam.svg?branch=master)](https://travis-ci.org/TomWright/queryparam)
[![codecov](https://codecov.io/gh/TomWright/queryparam/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/queryparam)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/queryparam)](https://goreportcard.com/report/github.com/TomWright/queryparam)
[![Documentation](https://godoc.org/github.com/TomWright/queryparam?status.svg)](https://godoc.org/github.com/TomWright/queryparam)

Stop accessing query strings and repeatedly parsing them into your preferred values - `queryparam` can do that for you! 

## Installation

```
go get -u github.com/tomwright/queryparam
```

```
import github.com/tomwright/queryparam/v4
```

## Examples

For examples see [godoc examples](https://godoc.org/github.com/TomWright/queryparam#example-Parse).

## Quickstart

Transform your http handlers from this...
```
func searchUsersHandler(r *http.Request, rw http.ResponseWriter) {
	values := r.URL.Query()

	userIDs := make([]string, 0)
	if userIDsStr := values.Get("id"); userIDsStr != "" {
		userIDs = strings.Split(userIDsStr, ",")
	}
	teamIDs := make([]string, 0)
	if teamIDsStr := values.Get("team-id"); teamIDsStr != "" {
		teamIDs = strings.Split(teamIDsStr, ",")
	}
	mustBeActive := false
	switch strings.ToLower(values.Get("must-be-active")) {
	case "true", "yes", "y":
		mustBeActive = true
	case "":
		break
	default:
		// unhandled bool value... handle as 400 or ignore to default as false
	}
	createdAfter := time.Time{}
	if createdAfterStr := values.Get("must-be-active"); createdAfterStr != "" {
		var err error
		createdAfter, err = time.Parse(time.RFC3339, createdAfterStr)
		if err != nil {
			// bad time value
		}
	}

	users, err := searchUsers(userIDs, teamIDs, mustBeActive, createdAfter)

	// handle users and err...
}
```

To this...
```
func searchUsersHandler(r *http.Request, rw http.ResponseWriter) {
	req := struct {
		UserIDs      []string  `queryparam:"id"`
		TeamIDs      []string  `queryparam:"team-id"`
		MustBeActive bool      `queryparam:"must-be-active"`
		CreatedAfter time.Time `queryparam:"created-after"`
	}{}

	err := queryparam.Parse(r.URL.Query(), &req)
	switch err {
	case nil:
		break
	case queryparam.ErrInvalidBoolValue: // only necessary if the request contains a bool value
		// they have entered a non-bool value.
		// this can be handled this as a 400 or ignored to default to false.
		return
	default:
		// something went wrong when parsing a value.
		// send a 500.
		return
	}

	users, err := searchUsers(req.UserIDs, req.TeamIDs, req.MustBeActive, req.CreatedAfter)

	// handle users and err...
}
```

## Types

By default `queryparam` can parse the following types.

- `string`
- `[]string`
- `int`
- `int32`
- `int64`
- `float32`
- `float64`
- `bool`
- `time.Time`
- `queryparam.Present`

### Custom Types

You can add custom type parsers and setters with the following:

```
// your custom type.
type MyCustomStringType string

// add a value parser for the custom type.
queryparam.DefaultParser.ValueParsers[reflect.TypeOf(MyCustomStringType(""))] = func(value string, _ string) (reflect.Value, error) {
    return reflect.ValueOf(MyCustomStringType(value)), nil
}
```

You can override the default value parsers in a similar manner...
```
queryparam.DefaultParser.ValueParsers[reflect.TypeOf("")] = func(value string, _ string) (reflect.Value, error) {
    // my custom string parser
    return reflect.ValueOf(value), nil
}
```
