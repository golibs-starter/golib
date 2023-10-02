package golib

import (
	"github.com/golibs-starter/golib/web/listener"
	"github.com/golibs-starter/golib/web/properties"
	"go.uber.org/fx"
)

func HttpRequestLogOpt() fx.Option {
	return fx.Options(
		ProvideEventListener(listener.NewRequestCompletedLogListener),
		ProvideProps(properties.NewHttpRequestLogProperties),
	)
}
