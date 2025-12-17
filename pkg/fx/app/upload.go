package app

import (
	"github.com/alanshaw/1up-service/pkg/fx/router"
	"github.com/alanshaw/1up-service/pkg/fx/upload/ucan"
	"go.uber.org/fx"
)

var UploadModule = fx.Module("upload",
	router.Module,
	ucan.Module,
)
