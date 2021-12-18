# Golib

Common core for Golang project.

### ⚠️ **Notice** ⚠️

Our modules are in private repo, so you need to config something bellow before start your develop.

#### 1. Setup `GOPRIVATE`

Run the command `export GOPRIVATE="gitlab.id.vin"` to add `gitlab.id.vin` as private repo.

For future usage, you might add above command to `.bashrc` or `.zshrc`.

#### 2. Add credentials to private host

Run the following command line to load `https://gitlab.id.vin/` using SSH:

```shell
git config --global url."git@gitlab.id.vin:".insteadOf "https://gitlab.id.vin/"
```

Or with access token in URL:

```shell
git config --global url."https://oath2:{your_access_token}@gitlab.id.vin/".insteadOf https://gitlab.id.vin/
```

### Setup instruction

Both `go get` and `go mod` are supported.

```shell
go get gitlab.id.vin/vincart/golib
```

### Usage

Using `fx.Option` to include dependencies for injection.

```go
package main

import (
	"context"
	"gitlab.id.vin/vincart/golib"
	"gitlab.id.vin/vincart/golib/actuator"
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/client"
	"go.uber.org/fx"
	"time"
)

func main() {
	_ = []fx.Option{
		// Required
		golib.AppOpt(),

		// Required
		golib.PropertiesOpt(),

		// When you want to use default logging strategy.
		golib.LoggingOpt(),

		// When you want to enable event publisher
		golib.EventOpt(),

		// When you want to enable actuator endpoints.
		// By default, we provide HealthService and InfoService.
		golib.ActuatorEndpointOpt(),
		// When you want to provide build info to above InfoService.
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		// When you want to provide custom health checker and informer
		golib.ProvideHealthChecker(NewSampleHealthChecker),
		golib.ProvideInformer(NewSampleInformer),

		// When you want to enable http client auto config with contextual client by default
		golib.HttpClientOpt(),

		// When you want to tell GoLib to load your properties.
		golib.ProvideProps(NewSampleProperties),

		// When you want to declare a service
		fx.Provide(NewSampleService),

		// When you want to register your event listener.
		golib.ProvideEventListener(NewSampleListener),
	}
}

// ==================================================
// ============ Example for build info ==============
// ==================================================
// These infos might inject in the build time
var (
	Version    = "1.0"
	CommitHash = "49d52932"
	BuildTime  = "2021/11/30 20:08:12"
)

// ==================================================
// ===== Example for register a health component ====
// ==================================================
// In this example, endpoint /actuator/health with return
//{
//  "meta": {
//    "code": 200,
//    "message": "Server is up"
//  },
//  "data": {
//    "status": "UP",
//    "components": {
//      "sample": {
//        "status": "UP"
//      }
//    }
//  }
//}

func NewSampleHealthChecker() actuator.HealthChecker {
	return &SampleHealthChecker{}
}

type SampleHealthChecker struct {
}

func (h SampleHealthChecker) Component() string {
	return "sample"
}

func (h SampleHealthChecker) Check(ctx context.Context) actuator.StatusDetails {
	// Do something to check status and
	// return actuator.StatusUp or actuator.StatusDown
	return actuator.StatusDetails{Status: actuator.StatusUp}
}

// ==================================================
// ========= Example for register a informer ========
// ==================================================
// In this example, endpoint /actuator/info with return
//{
//  "meta": {
//    "code": 200,
//    "message": "Successful"
//  },
//  "data": {
//    "service_name": "Sample Service",
//    "info": {
//      "sample": {
//        "key1": "val1"
//      }
//    }
//  }
//}

func NewSampleInformer() actuator.Informer {
	return &SampleInformer{}
}

type SampleInformer struct {
}

func (s SampleInformer) Key() string {
	return "sample"
}

func (s SampleInformer) Value() interface{} {
	return map[string]interface{}{"key1": "val1"}
}

// ==================================================
// ======== Example about declare properties ========
// ==================================================

func NewSampleProperties(loader config.Loader) (*SampleProperties, error) {
	props := SampleProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type SampleProperties struct {
	// We use github.com/go-playground/validator to validate properties
	Field1 string `validate:"required"`

	// We use https://github.com/zenthangplus/defaults to set default for properties
	Field2 int `default:"10"`

	// We use github.com/mitchellh/mapstructure to bind config to properties
	Field3 []time.Duration `mapstructure:"field3_new_name"`
}

// Prefix
// Defines the properties prefix
func (s SampleProperties) Prefix() string {
	return "app.sample"
}

// ==================================================
// ======== Example about inject dependencies =======
// ==================================================

// NewSampleService In this case Contextual Http Client is required
func NewSampleService(httpClient client.ContextualHttpClient) *SampleService {
	return &SampleService{httpClient: httpClient}
}

type SampleService struct {
	httpClient client.ContextualHttpClient
}

// ==================================================
// === Example about declare listener (subscriber) ==
// ==================================================

func NewSampleListener() pubsub.Subscriber {
	return &SampleListener{}
}

type SampleListener struct {
}

func (s SampleListener) Supports(e pubsub.Event) bool {
	_, ok := e.(*event.RequestCompletedEvent)
	return ok
}

func (s SampleListener) Handle(e pubsub.Event) {
	// Handle when receive event
}

```

> More complex examples in [golib-sample](../golib-sample)

### Configuration

#### 1. Environment variables

| Var                         | Default    | Description                                                                                                                                                                                                                                |
|-----------------------------|------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `APP_PROFILES` or `APP_ENV` | None       | Defines the list of active profiles, separate by comma. By default, `default` profile is always load even this env configured. Example: when `APP_PROFILES=internal,uat` then both `default` `internal` and `uat` will be loaded by order. |
| `APP_CONFIG_PATHS`          | `./config` | Defines the location of config directory, when the application is started, it will scan profiles in this path.                                                                                                                             |
| `APP_CONFIG_FORMAT`         | `yaml`     | Defines the format of config file. Currently we only support Yaml format (both `yaml` `yml` are accepted).                                                                                                                                 |

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
