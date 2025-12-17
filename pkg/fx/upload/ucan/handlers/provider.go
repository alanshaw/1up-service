package handlers

import (
	"go.uber.org/fx"

	"github.com/alanshaw/1up-service/pkg/service/upload/ucan"
)

var Module = fx.Module("storage/ucan/handlers",
	fx.Provide(
		fx.Annotate(
			ucan.NewDebugEchoHandler,
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
)
