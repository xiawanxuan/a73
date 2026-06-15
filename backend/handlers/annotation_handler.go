package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ocean-survey-backend/services"
	"ocean-survey-backend/utils"
)

type AnnotationHandler struct {
	Service *services.AnnotationService
}

func NewAnnotationHandler(service *services.AnnotationService) *AnnotationHandler {
	return &AnnotationHandler{Service: service}
}

func getUserID(c *gin.Context) uuid.UUID {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func (h *AnnotationHandler) ListByPointCloud(c *gin.Context) {
	pcIDStr := c.Param("pcID")
	pcID, err := uuid.Parse(pcIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid point cloud id")
		return
	}

	list, err := h.Service.ListByPointCloud(pcID)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, list)
}

func (h *AnnotationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	ann, err := h.Service.GetByID(id)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, ann)
}

func (h *AnnotationHandler) Create(c *gin.Context) {
	pcIDStr := c.Param("pcID")
	pcID, err := uuid.Parse(pcIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid point cloud id")
		return
	}

	var dto services.CreateAnnotationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)

	ann, err := h.Service.ValidateAndCreate(pcID, &dto, userID)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, ann)
}

func (h *AnnotationHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	var dto services.UpdateAnnotationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)

	ann, err := h.Service.ValidateAndUpdate(id, &dto, userID)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, ann)
}

func (h *AnnotationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	userID := getUserID(c)

	if err := h.Service.Delete(id, userID); err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *AnnotationHandler) Rollback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	versionStr := c.Query("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		utils.Fail(c, http.StatusBadRequest, "invalid version")
		return
	}

	userID := getUserID(c)

	ann, err := h.Service.Rollback(id, version, userID)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, ann)
}

func (h *AnnotationHandler) ListSnapshots(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	snapshots, err := h.Service.ListSnapshots(id)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, snapshots)
}
