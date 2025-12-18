package handlers

import (
	"go.uber.org/fx"

	"github.com/alanshaw/1up-service/pkg/service/upload/ucan"
)

var Module = fx.Module("upload/ucan/handlers",
	fx.Provide(
		fx.Annotate(
			ucan.NewDebugEchoHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			ucan.NewProviderDeregisterHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			ucan.NewProviderListHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			ucan.NewProviderRegisterHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			ucan.NewProviderWeightSetHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			ucan.NewBlobAddHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
)
