package app

import (
	"github.com/alanshaw/1up-service/pkg/fx/upload/ucan"
	"github.com/alanshaw/1up-service/pkg/fx/upload/ucan/handlers"
	"go.uber.org/fx"
)

var UploadModule = fx.Module("upload",
	ucan.Module,
	handlers.Module,
)
