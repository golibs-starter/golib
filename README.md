# Golib

Common core for Golang project.

### ⚠️ **Notice** ⚠️

Our modules are in private repo, so you need to config something bellow before start your develop.

#### 1. Setup `GOPRIVATE`

Run the command `export GOPRIVATE="gitlab.com"` to add `gitlab.com` as private repo.

For future usage, you might add above command to `.bashrc` or `.zshrc`.

#### 2. Add credentials to private host

Run the following command line to load `https://gitlab.com/` using SSH:

```shell
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```

Or with access token in URL:

```shell
git config --global url."https://oath2:{your_access_token}@gitlab.com/".insteadOf https://gitlab.com/
```

### Setup instruction

Both `go get` and `go mod` are supported.

```shell
go get gitlab.com/golibs-starter/golib
```

### Usage

Using `fx.Option` to include dependencies for injection.

See examples:

- [Bootstrap your application](./example/bootstrap.go)
- [Declare a properties](./example/sample_properties.go)
- [Declare an event](./example/sample_event.go)
- [Declare a service](./example/sample_service.go)
- [Declare a listener (subscriber)](./example/sample_listener.go)
- [Provide build info](./example/samle_build_info.go)
- [Register an informer](./example/sample_informer.go)
- [Register a health checker](./example/sample_health_checker.go)

> More complex examples in [golib-sample](https://gitlab.com/golibs-starter/golib-sample)

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

    # Configuration available for EventOpt()
    event:
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
```
