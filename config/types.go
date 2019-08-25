package config

type AppConfig struct {
	FileConfig   FileConfig
	MemoryConfig MemoryConfig
}

type MemoryConfig struct {
	AppPath        string // app根目录
	ConfigFileName string //配置文件名
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
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Database string `json:"Database"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type LogConfig struct {
	LogDir       string `json:"LogDir"`
	BackupNumber int    `json:"BackupNumber"`
	MaxSize      int    `json:"MaxSize"`
}

type WeiXinConfig struct {
	AppId  string `json:"AppId"`
	Secret string `json:"Secret"`
	Token  string `json:"Token"`
}
