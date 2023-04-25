# Meater!Go
Is a simple API client to interact with MEATER Cloud REST API service (https://public-api.cloud.meater.com) via HTTP, written in Go

## News
  * v1.0.0 First public release on April 245, 2023

## Features
  * Simple and chainable methods for settings and request

## Installation
```bash
# Go Modules
go get github.com/alex-savin/go-meater
```

## Usage
The following samples will assist you to become as comfortable as possible with Meater!Go library.
```go
// Import Meater!Go into your code and refer it as `meater`.
import "github.com/alex-savin/go-meater"
```

#### Create a new Meater connection
```go
// Create a Meater API Client
client, _ := meater.New(
	meater.Username("username"),
	meater.Password("password"),
	meater.LogLevel("log"),
)
```

#### Looping over active probes
```go
if len(client.GetProbes()) > 0 {
	for i, probe := range client.GetProbes() {
		fmt.Printf("PROBE #%d: %+v\n", i+1, probe)
		fmt.Printf("COOK  #%d: %+v\n", i+1, probe.Cook)
	}
} else {
    fmt.Print("No active probes are detected\n")
}
```