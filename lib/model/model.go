package model

// エラーレスポンス用の構造体
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
