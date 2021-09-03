package v1

import (
	"github.com/cookienyancloud/avito-backend-test/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.FinanceService
}

func NewHandler(services *service.FinanceService) *Handler {
	return &Handler{
		services:     services,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initFinanceRoutes(v1)
	}
}