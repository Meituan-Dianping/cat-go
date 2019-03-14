package message

type Header struct {
	Domain   string
	Hostname string
	Ip       string

	MessageId       string
	ParentMessageId string
	RootMessageId   string
}
