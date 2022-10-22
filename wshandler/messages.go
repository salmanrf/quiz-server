package wshandler

type Message[MsgType interface{}] struct {
	Event string `json:"event"`
	Data MsgType `json:"data"`
	sender *Member 
}