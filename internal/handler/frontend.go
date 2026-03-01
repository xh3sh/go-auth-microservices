package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleHome СЂРµРЅРґРµСЂРёС‚ РіР»Р°РІРЅСѓСЋ СЃС‚СЂР°РЅРёС†Сѓ РїСЂРёР»РѕР¶РµРЅРёСЏ
func (h *Handler) HandleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}
