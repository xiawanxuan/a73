package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ocean-survey-backend/services"
	"ocean-survey-backend/utils"
)

type PointCloudHandler struct {
	Service *services.PointCloudService
}

func NewPointCloudHandler(service *services.PointCloudService) *PointCloudHandler {
	return &PointCloudHandler{Service: service}
}

func (h *PointCloudHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	list, total, err := h.Service.List(page, pageSize)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PageResult(c, list, total, page, pageSize)
}

func (h *PointCloudHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	pc, err := h.Service.GetByID(id)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, pc)
}

func (h *PointCloudHandler) Create(c *gin.Context) {
	var dto services.CreatePointCloudDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	pc, err := h.Service.Create(&dto)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, pc)
}

func (h *PointCloudHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	var dto services.UpdatePointCloudDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	pc, err := h.Service.Update(id, &dto)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, pc)
}

func (h *PointCloudHandler) Delete(c *gin.Context) {
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
