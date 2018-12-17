# OpenTracing support for AWS SDK in Go

The `otaws` package makes it easy to add OpenTracing support for AWS SDK in Go.

## Installation

```
go get github.com/opentracing-contrib/go-aws
```

## Documentation

See the basic usage examples below and the [package documentation on
godoc.org](https://godoc.org/github.com/opentracing-contrib/go-aws).

## Usage

```go
// You must have some sort of OpenTracing Tracer instance on hand
var tracer opentracing.Tracer = ...

// Set Tracer as global 
opentracing.SetGlobalTracer(tracer)

// Create AWS Session
sess := session.NewSession(...)

// Create AWS service client e.g. DynamoDB client
dbCient := dynamodb.New(sess)

// Add OpenTracing handlers:
AddOTHandlers(dbClient.Client)

// Call AWS client
result, err := dbClient.ListTables(&dynamodb.ListTablesInput{})

```

## License

[Apache 2.0 License](./LICENSE).
