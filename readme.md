# Traefik Request Timestamps Middleware

A Traefik plugin that adds request and response timestamp headers to HTTP responses, helping you track request processing times and latency.

This plugin injects both a REQUEST-TIMESTAMP and a RESPONSE-TIMESTAMP header into each HTTP response. The REQUEST-TIMESTAMP reflects when the request first reached Traefik, while the RESPONSE-TIMESTAMP marks when the response is sent back to the client. Both header names and the date format are fully configurable, supporting any Go time format. The plugin is lightweight, fast, and has no dependencies.

## Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `requestHeaderName` | Name of the request timestamp header | `"REQUEST-TIMESTAMP"` |
| `responseHeaderName` | Name of the response timestamp header | `"RESPONSE-TIMESTAMP"` |
| `dateFormat` | Go time format string | `"2006-01-02T15:04:05.000Z"` |

### Static Configuration (traefik.yml)

```yaml
experimental:
  plugins:
    timestampheaders:
      moduleName: "github.com/feedzai/traefikRequestTimestamps"
      version: "v1.0.0"
```

### Dynamic Configuration (dynamic.yml)

#### Basic Usage (Default Settings)
```yaml
http:
  middlewares:
    timestamp-headers:
      plugin:
        timestampheaders: {}

  routers:
    api:
      rule: "Path(`/`)"
      service: backend-service
      middlewares:
        - timestamp-headers
```

#### Custom Configuration
```yaml
http:
  middlewares:
    custom-timestamps:
      plugin:
        timestampheaders:
          requestHeaderName: "X-Request-Time"
          responseHeaderName: "X-Response-Time"
          dateFormat: "2006-01-02 15:04:05 UTC"
```

## Requirements

- Traefik v3.0+
- Go 1.19+ (for development)
