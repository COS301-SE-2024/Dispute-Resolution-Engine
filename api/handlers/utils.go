package handlers

import (
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Utility) GetCountries(c *gin.Context) {
	var countries []models.Country
	err := h.DB.Find(&countries).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Data: countries})
}
