package logic

//广播器

type broadcaster struct {
}

var Broadcaster = &broadcaster{}

func (b *broadcaster) Broadcast(msg *Message) {
	return
}
