package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ocean-survey-backend/services"
	"ocean-survey-backend/utils"
	"ocean-survey-backend/websocket"
)

type LabelHandler struct {
	Service *services.TerrainLabelService
	Hub     *websocket.Hub
}

func NewLabelHandler(service *services.TerrainLabelService, hub *websocket.Hub) *LabelHandler {
	return &LabelHandler{Service: service, Hub: hub}
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil, false
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false
	}
	return uid, true
}

func (h *LabelHandler) List(c *gin.Context) {
	var userID *uuid.UUID
	if uid, ok := getUserID(c); ok {
		userID = &uid
	}

	list, err := h.Service.List(userID)
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

func (h *LabelHandler) broadcastLabel(labelMsgType string, labelID string, userID string, payload any) {
	if h.Hub == nil {
		return
	}
	labelData := websocket.LabelData{
		LabelID: labelID,
		UserID:  userID,
		Payload: payload,
	}
	dataBytes, _ := json.Marshal(labelData)
	msg := &websocket.WSMessage{
		Type:      labelMsgType,
		Data:      dataBytes,
		Timestamp: time.Now().UnixMilli(),
	}
	h.Hub.BroadcastAll <- msg
}

func (h *LabelHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "X-User-ID header is required")
		return
	}

	var dto services.CreateTerrainLabelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	label, err := h.Service.Create(userID, &dto)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.broadcastLabel(websocket.MsgTypeLabelCreate, label.ID.String(), userID.String(), label)

	utils.Success(c, label)
}

func (h *LabelHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "X-User-ID header is required")
		return
	}

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

	label, err := h.Service.Update(id, userID, &dto)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	h.broadcastLabel(websocket.MsgTypeLabelUpdate, label.ID.String(), userID.String(), label)

	utils.Success(c, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "X-User-ID header is required")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Service.Delete(id, userID); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	h.broadcastLabel(websocket.MsgTypeLabelDelete, id.String(), userID.String(), nil)

	utils.Success(c, nil)
}
