package cmessage

type CMessageTemplate interface {
	FillParams(params map[string]string) (content string)
}

type CMessage interface {
	Initialize(params map[string]string)
	SetTemplate(template CMessageTemplate)
	Send(from string, to []string, params map[string]string) error
}

var cmessageMap = map[string]func() CMessage{}

func register(name string, f func() CMessage) {
	cmessageMap[name] = f
}

func New(name string) CMessage {
	f, ok := cmessageMap[name]
	if !ok {
		return cmessageMap["unknown"]()
	}
	return f()
}
