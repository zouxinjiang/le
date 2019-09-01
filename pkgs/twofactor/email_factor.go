package twofactor

type EmailFactor struct {
}

func init() {
	register("email", newEmailFactor)
}

func newEmailFactor() TwoFactor {
	return &EmailFactor{}
}

func (e EmailFactor) Do(params map[string]string) (string, error) {
	panic("implement me")
}
