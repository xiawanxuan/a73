package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ocean-survey-backend/services"
	"ocean-survey-backend/utils"
)

type LabelHandler struct {
	Service *services.TerrainLabelService
}

func NewLabelHandler(service *services.TerrainLabelService) *LabelHandler {
	return &LabelHandler{Service: service}
}

func (h *LabelHandler) List(c *gin.Context) {
	list, err := h.Service.List()
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, list)
}

func (h *LabelHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	label, err := h.Service.GetByID(id)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, label)
}

func (h *LabelHandler) Create(c *gin.Context) {
	var dto services.CreateTerrainLabelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	label, err := h.Service.Create(&dto)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, label)
}

func (h *LabelHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	var dto services.UpdateTerrainLabelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	label, err := h.Service.Update(id, &dto)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Service.Delete(id); err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, nil)
}
