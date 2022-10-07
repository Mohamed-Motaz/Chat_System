package Server

type AddEntityRes struct {
	Number int32 `json:"number"`
}

type AddMessageReq struct {
	Body string `json:"body"`
}
type AddMessageRes struct {
	Number int32 `json:"number"`
}

type UpdateMessageReq struct {
	Body string `json:"body"`
}

type UpdateMessageRes struct {
	Number int32  `json:"number"`
	Body   string `json:"body"`
}

type SearchForMessageReq struct {
	Body string `json:"body"`
}

type MessageRes struct {
	Number int32  `json:"number"`
	Body   string `json:"body"`
}
type SearchForMessageRes struct {
	Messages []MessageRes `json:"Messages"`
}
