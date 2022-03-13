package golib

import (
	"gitlab.com/golibs-starter/golib/web/listener"
	"gitlab.com/golibs-starter/golib/web/properties"
	"go.uber.org/fx"
)

func HttpRequestLogOpt() fx.Option {
	return fx.Options(
		ProvideEventListener(listener.NewRequestCompletedLogListener),
		ProvideProps(properties.NewHttpRequestLogProperties),
	)
}
