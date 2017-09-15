# consul

## TL;DR
Fetch IP:port of a service from the local Consul datacentre. Load balancing is
performed if there is more than one service running under the same name.

## Description
Package `consul` currently contains a single function (`ServiceHostPort`) which
takes a service name and uses Consul's DNS server to fetch its SRV record if
one exists. An SRV record contains both IP address and port of the service.

If there is more than one service registered with the same name, the response
for such DNS request will contain multiple IP:port pairs and one will be chosen
at random.

## Installation
`consul` can be imported with your vendoring tool:
```
go get github.com/fresh8/consul
```

Once imported into your project start referencing it by adding the line
```go
import "github.com/fresh8/consul"
```

## Examples
```go
addr, err := consul.ServiceHostPort("redis")
if err == nil {
  fmt.Printf("Redis is at %s", addr) // → Redis is at 172.16.0.1:32280
}
```

## Tests
In the package root directory run `go test` to run all tests.
