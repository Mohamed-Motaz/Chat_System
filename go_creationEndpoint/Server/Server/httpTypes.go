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
