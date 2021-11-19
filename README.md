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
options := []fx.Option{
    golib.AppOpt(),//required
    golib.PropertiesOpt(),//required

    // When you want to use default logging strategy.
    golib.LoggingOpt(),

    // When you want to enable event publisher
    golib.EventOpt(),

    // When you want to enable actuator endpoints.
	// By default, we provide HealthService and InfoService.
    golib.ActuatorEndpointOpt(),
    // When you want to provide build info to above info service.
    golib.BuildInfoOpt(Version, CommitHash, BuildTime),

    //When you want to enable http client auto config with contextual client by default
    golib.HttpClientOpt(),
    //When you want to provide an additional wrapper to easy to control http client's security.
    golibsec.SecuredHttpClientOpt(),

    // When you want to tell GoLib to load your properties.
    golib.ProvideProps(properties.NewCustomProperties),

    // When you want to register your event listener.
	golib.ProvideEventListener(listener.NewCustomListener),
}
```
