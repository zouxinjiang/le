package gongzhonghao

type ApiAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ApiError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type CustomMenuButtonType string

const (
	BtnType_Click CustomMenuButtonType = "click"
	BtnType_View  CustomMenuButtonType = "view"
	BtnType_Media CustomMenuButtonType = "media_id"
)

type CustomMenuButton struct {
	Type      *CustomMenuButtonType `json:"type,omitempty"`
	Name      string                `json:"name,omitempty"`
	Key       *string               `json:"key,omitempty"`
	Url       *string               `json:"url,omitempty"`
	MediaId   *string               `json:"media_id,omitempty"`
	SubButton []CustomMenuButton    `json:"sub_button,omitempty"`
}

type Industry struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
}
type ApiIndustry struct {
	PrimaryIndustry   Industry `json:"primary_industry"`
	SecondaryIndustry Industry `json:"secondary_industry"`
}

type MessageTemplate struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

type TemplateMessageDataItem struct {
	Keyword string
	Value   string
	Color   string
}
