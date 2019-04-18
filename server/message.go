package main

type Message interface {
	GetSenderId() (int, bool)
	GetMessageContent() (string, bool)
}

type IdMessage struct {
	id      int
	content string
}

func (m IdMessage) GetSenderId() (int, bool) {
	return m.id, true
}

func (m IdMessage) GetContent() (string, bool) {
	return m.content, true
}

type PlainMessage struct {
	content string
}

func (m PlainMessage) GetSenderId() (int, bool) {
	return 0, false
}

func (m PlainMessage) GetContent() (string, bool) {
	return m.content, true
}

type NilMessage struct {
	id int
}

func (m NilMessage) GetSenderId() (int, bool) {
	return m.id, true
}

func (m NilMessage) GetContent() (string, bool) {
	return "", false
}
