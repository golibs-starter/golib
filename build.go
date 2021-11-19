package golib

import (
	"gitlab.id.vin/vincart/golib/build"
	"go.uber.org/fx"
)

func BuildInfoOpt(version string, commitHash string, time string) fx.Option {
	return fx.Options(
		fx.Supply(build.Version(version)),
		fx.Supply(build.CommitHash(commitHash)),
		fx.Supply(build.Time(time)),
		fx.Provide(fx.Annotated{
			Group:  "actuator_informer",
			Target: build.NewInformer,
		}),
	)
}
