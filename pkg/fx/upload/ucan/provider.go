package ucan

import (
	echofx "github.com/alanshaw/1up-service/pkg/fx/echo"
	"github.com/alanshaw/1up-service/pkg/fx/upload/ucan/handlers"
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/ucantone/principal"
	"github.com/alanshaw/ucantone/server"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var log = logging.Logger("fx/upload/ucan")

type Server struct {
	ucanServer *server.HTTPServer
}

var Module = fx.Module("upload/ucan/server",
	fx.Provide(
		NewServer,
		fx.Annotate(
			NewServer,
			fx.As(new(echofx.RouteRegistrar)),
			fx.ResultTags(`group:"route_registrar"`),
		),
	),
	handlers.Module,
)

type Params struct {
	fx.In
	ID       principal.Signer
	Handlers []*service.Handler  `group:"ucan_handlers"`
	Options  []server.HTTPOption `group:"ucan_options"`
}

func NewServer(p Params) (*Server, error) {
	ucanSvr := server.NewHTTP(p.ID, p.Options...)
	log.Infof("Registering %d UCAN handlers", len(p.Handlers))
	for _, h := range p.Handlers {
		log.Infof("%q", h.Capability.Command())
		ucanSvr.Handle(h.Capability, h.Handler)
	}
	return &Server{ucanSvr}, nil
}

func (s *Server) RegisterRoutes(e *echo.Echo) {
	e.POST("/", echo.WrapHandler(s.ucanServer))
}
