package wshandler

var MemberEvents = make(map[string]*func(*Member))

func StartQuiz(m *Member, payload Message[interface{}]) {
	m.room.broadcast <- []byte(payload.Event)
}