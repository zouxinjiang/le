package web

type ApiAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

type ApiUserInfo struct {
	OpenId     string `json:"openid"`
	NickName   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

type ApiError struct {
	ErrorCode int64  `json:"errorcode"`
	ErrMsg    string `json:"errmsg"`
}
