package twofactor

type TwoFactor interface {
	Do(params map[string]string) (string, error)
}

type TwoFactorType = string

type TwoFactorInstanceFunc func() TwoFactor

var tfmap = map[TwoFactorType]TwoFactorInstanceFunc{}

func register(name TwoFactorType, value func() TwoFactor) {
	tfmap[name] = value
}

func New(name TwoFactorType) TwoFactor {
	f, ok := tfmap[name]
	if ok {
		return tfmap["unknown"]()
	}
	return f()
}
