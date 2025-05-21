# Gometrix - Metrics Client Library for Go

Gometrix is a flexible metrics client library for Go applications that supports multiple backends including StatsD, logging to file, and a dummy client for testing.

## Features

- Multiple client implementations:
  - **StatsD client** - Sends metrics to a StatsD server
  - **Logging client** - Writes metrics to log files with aggregation
  - **Dummy client** - No-op implementation for testing
- Common metrics interface:
  - Counters
  - Gauges
  - Timings
  - Increment/Decrement operations
- Tag support for all metric types
- Thread-safe implementations

## Installation

```bash
go get github.com/allread-mlt/gometrix
```

## Usage

### Basic Configuration
First, import the package:

```Go
import "github.com/allread-mlt/gometrix"
```

### Creating a Client
```Go
config := &gometrix.MetricsClientConfig{
    ClientType: gometrix.metricsTypeStatsd, // or metricsTypeLogging/metricsTypeDummy
    ClientData: map[string]any{
        // Configuration specific to your client type
    },
}

client, err := gometrix.NewMetricsClient(config)
if err != nil {
    log.Fatalf("Failed to create metrics client: %v", err)
}
defer client.Stop()
```

### Client-Specific Configurations

#### StatsD Client
```Go
clientData := StatsdMetricsData{
    ServerHost:   "statsd.example.com",
    ServerPort:   8125,
    Prefix: "myapp.",
}
```

#### Logging Client
```Go
clientData := LoggingMetricsData{
    Timeout:       60, // seconds between log outputs
    LogFilePath: "/var/log/myapp",
    MaxFiles:     3,
    MaxFileSize: 10, // MB
}
```

### Sending Metrics

```Go
// Increment a counter
client.Increment("user.logins", 1, map[string]any{"source": "web"})

// Record a gauge value
client.Gauge("memory.usage", 85.3, nil)

// Measure duration
start := time.Now()
// ... some operation ...
client.Timing("db.query_time", time.Since(start), map[string]any{"query": "get_user"})
```

## Configuration Examples using YAML

### statsD client

```YAML
metrics:
  type: statsd
  data:
    host: "localhost"
    port: 8125
    prefix: "app.prod."
```

### Logging client

```YAML
metrics:
  type: logging
  data:
    timeout: 60 # in seconds
    log_file_path: "metrics"
    max_files: 3
    max_file_size: 10 # in mb
```


### Dummy client

```YAML
metrics:
  type: dummy
```
