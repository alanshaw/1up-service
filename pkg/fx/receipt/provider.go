package receipt

import (
	echofx "github.com/alanshaw/1up-service/pkg/fx/echo"
	"github.com/alanshaw/1up-service/pkg/store/token"
	"github.com/alanshaw/libracha/receipt"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("receipt_handler",
	fx.Provide(
		fx.Annotate(
			NewReceiptHandler,
			fx.As(new(echofx.RouteRegistrar)),
			fx.ResultTags(`group:"route_registrar"`),
		),
	),
)

var _ echofx.RouteRegistrar = (*Handler)(nil)

type Handler struct {
	tokens token.Store
}

func NewReceiptHandler(tokens token.Store) *Handler {
	return &Handler{tokens}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/receipt/:task", echo.WrapHandler(receipt.NewHandler(h.tokens)))
}
