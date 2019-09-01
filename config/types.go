package config

type AppConfig struct {
	FileConfig   FileConfig
	MemoryConfig MemoryConfig
}

type MemoryConfig struct {
	AppPath        string // app根目录
	ConfigFileName string //配置文件名
	EncryptKey     string // 加密/混淆密码
	LoginKey       string // 验证登陆信息的key。 登陆session的key和token认证的key
}

type FileConfig struct {
	WebConfig      WebConfig      `json:"WebConfig"`
	DatabaseConfig DatabaseConfig `json:"DatabaseConfig"`
	LogConfig      LogConfig      `json:"LogConfig"`
	WeiXinConfig   WeiXinConfig   `json:"WeiXinConfig"`
}

type WebConfig struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type DatabaseConfig struct {
	Debug    bool   `json:"Debug"`
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Database string `json:"Database"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type LogConfig struct {
	ShowLevel    uint64 `json:"ShowLevel"`
	LogDir       string `json:"LogDir"`
	BackupNumber int    `json:"BackupNumber"`
	MaxSize      int    `json:"MaxSize"`
}

type WeiXinConfig struct {
	AppId  string `json:"AppId"`
	Secret string `json:"Secret"`
	Token  string `json:"Token"`
}
