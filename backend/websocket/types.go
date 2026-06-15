package websocket

import "encoding/json"

const (
	MsgTypeAnnotationCreate   = "annotation_create"
	MsgTypeAnnotationUpdate   = "annotation_update"
	MsgTypeAnnotationDelete   = "annotation_delete"
	MsgTypeAnnotationRollback = "annotation_rollback"
	MsgTypeLabelCreate        = "label_create"
	MsgTypeLabelUpdate        = "label_update"
	MsgTypeLabelDelete        = "label_delete"
	MsgTypeUserJoin           = "user_join"
	MsgTypeUserLeave          = "user_leave"
	MsgTypeOnlineUsers        = "online_users"
	MsgTypeCursorMove         = "cursor_move"
	MsgTypeDraftSync          = "draft_sync"
	MsgTypePing               = "ping"
	MsgTypePong               = "pong"
)

type WSMessage struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

type UserJoinData struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Color    string `json:"color"`
}

type UserLeaveData struct {
	UserID string `json:"userId"`
}

type AnnotationData struct {
	AnnotationID string `json:"annotationId"`
	PointCloudID string `json:"pointCloudId"`
	Payload      any    `json:"payload"`
}

type CursorData struct {
	UserID       string  `json:"userId"`
	PointCloudID string  `json:"pointCloudId"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Z            float64 `json:"z"`
	Color        string  `json:"color,omitempty"`
}

type DraftData struct {
	UserID       string `json:"userId"`
	PointCloudID string `json:"pointCloudId"`
	Data         string `json:"data"`
}

type LabelData struct {
	LabelID string `json:"labelId"`
	UserID  string `json:"userId,omitempty"`
	Payload any    `json:"payload"`
}

type OnlineUsersData struct {
	Users []UserJoinData `json:"users"`
}
