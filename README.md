# Golib

> **Note**
> We are moving out from [Gitlab](https://gitlab.com/golibs-starter). All packages are now migrated to `github.com/golibs-starter/*`. Please consider updating.

[![run tests](https://github.com/golibs-starter/golib/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/golibs-starter/golib/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/golibs-starter/golib/graph/badge.svg?token=0CJ3MKCSVP)](https://codecov.io/gh/golibs-starter/golib)
[![goreportcard](https://goreportcard.com/badge/github.com/golibs-starter/golib)](https://goreportcard.com/report/github.com/golibs-starter/golib)

Common core for Golang project.

### Setup instruction

Both `go get` and `go mod` are supported.

```shell
go get github.com/golibs-starter/golib
```

### Usage

Using `fx.Option` to include dependencies for injection.

See some simple examples:

- [Bootstrap your application](./example/bootstrap.go)
- [Declare a properties](./example/sample_properties.go)
- [Declare an event](./example/sample_event.go)
- [Declare a service](./example/sample_service.go)
- [Declare a listener (subscriber)](./example/sample_listener.go)
- [Provide build info](./example/samle_build_info.go)
- [Register an informer](./example/sample_informer.go)
- [Register a health checker](./example/sample_health_checker.go)

> Full working examples in [golib-sample](https://github.com/golibs-starter/golib-sample):
> - [Public API Service (JWT Auth)](https://github.com/golibs-starter/golib-sample/-/tree/develop/src/public)
> - [Internal API Service (Basic Auth)](https://github.com/golibs-starter/golib-sample/-/tree/develop/src/internal)
> - [Worker Service](https://github.com/golibs-starter/golib-sample/-/tree/develop/src/worker)
> - [Migration Job](https://github.com/golibs-starter/golib-sample/-/tree/develop/src/migration)

### Configuration

#### 1. Environment variables

| Var                         | Default    | Description                                                                                                                                                                                                                                          |
|-----------------------------|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `APP_PROFILES` or `APP_ENV` | `local`    | Defines the list of active profiles, separate by comma. **By default, `default` profile is always load even this env configured**. <br/> Example: when `APP_PROFILES=internal,uat` then both `default` `internal` and `uat` will be loaded by order. |
| `APP_CONFIG_PATHS`          | `./config` | Defines the location of config directory, when the application is started, it will scan profiles in this path.                                                                                                                                       |
| `APP_CONFIG_FORMAT`         | `yaml`     | Defines the format of config file. Currently we only support Yaml format (both `yaml` `yml` are accepted).                                                                                                                                           |

Besides, all our configs can be overridden by environment variables. For example:

```yaml
store:
    name: Fruit store # Equivalent to STORE_NAME
    items:
        -   name: Apple # Equivalent to STORE_ITEMS_0_NAME
            price: 5 # Equivalent to STORE_ITEMS_0_PRICE
        -   name: Lemon # Equivalent to STORE_ITEMS_1_NAME
            price: 0.5 # Equivalent to STORE_ITEMS_1_PRICE
```

#### 2. Available configurations

```yaml
app:
    # Configuration available for AppOpt()
    name: Service Name # Specify application name. Default `unspecified`
    port: 8080 # Defines the running port. Default `8080`
    path: /service-base-path/ # Defines base path (context path). Default `/`

    # Configuration available for LoggingOpt()
    logging:
        development: false # Enable or disable development mode. Default `false`
        jsonOutputMode: true # Enable or disable json output. Default `true`
        logLevel: INFO # LogLevel is the minimum enabled logging level.

    # Configuration available for EventOpt()
    event:
        # Default event channel accept maximum 10 events,
        # other incoming events will be blocked until the channel has free space.
        channelSize: 10
        notLogPayloadForEvents:
            - OrderCreatedEvent
            - OrderUpdatedEvent

    # Configuration for HttpClientOpt()
    httpClient:
        timeout: 60s # Request timeout, in duration format. Default 60s
        maxIdleConns: 100 # Default 100
        maxIdleConnsPerHost: 10 # Default 10
        maxConnsPerHost: 100 # Default 100
        proxy:
            url: http://localhost:8080 # Proxy url
            appliedUris: # List of URIs, which will be requested under above proxy
                - https://foo.com/path/
                - https://bar.com/path/
    httpRequest:
        logging:
            disabled: false # Enable/disable all http request log
            predefinedDisabledUrls: # Shouldn't modify it, use disabledUrls instead
                - { urlPattern: "^/actuator/.*" } # By default, we disable all actuator requests
            disabledUrls: # Not log request for urls that matching method & url pattern
                - { method: "GET", urlPattern: "^/an-url-with-disabled-log/.*" }
                - { method: "POST", urlPattern: "^/another-url$" }

```
