package config

type Properties interface {
	Prefix() string
}

type PropertiesPreBinding interface {
	PreBinding() error
}

type PropertiesPostBinding interface {
	PostBinding() error
}
