package messenger

type Messenger interface {
	Send(message string) bool
	Receiver()(message string, ok bool)
}

type Telegram struct {

}

func (t *Telegram) Send(message string) bool {
	retuen true
}

func (t *Telegram) Receiver()(message string, ok bool) {
	return "", true
}
