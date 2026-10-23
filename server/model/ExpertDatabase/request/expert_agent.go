package request

// ExpertAgentChatReq 的人数参数为零时沿用会话条件；首次请求使用默认值。
type ExpertAgentChatReq struct {
	Message    string `json:"message"`
	SessionID  string `json:"sessionId"`
	Count      int    `json:"count"`
	MaxPerUnit int    `json:"maxPerUnit"`
}
