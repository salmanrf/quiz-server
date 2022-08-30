package wshandler

type Room struct {
	Code string `json:"code"`
	Quota int `json:"quota"`
	UsedQuota int `json:"used_quota"`
	members map[*Member]Member;
	broadcast chan []byte;
	register chan *Member;
	unregister chan *Member;
}