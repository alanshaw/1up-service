package router

import (
	"github.com/alanshaw/1up-service/pkg/service/router"
	"go.uber.org/fx"
)

var Module = fx.Module("router",
	fx.Provide(router.NewRouter),
)
