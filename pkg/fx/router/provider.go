package router

import (
	"github.com/alanshaw/1up-service/pkg/service/routing"
	"go.uber.org/fx"
)

var Module = fx.Module("router",
	fx.Provide(routing.NewRouter),
)
