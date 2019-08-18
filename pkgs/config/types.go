package config

type (
	AppConfig struct {
		Inited         bool           `json:"Inited"`
		ConfigDir      string         `json:"ConfigDir"`
		TemplateDir    string         `json:"TemplateDir"`
		WebDir         string         `json:"WebDir"`
		BinDir         string         `json:"BinDir"`
		DatabaseConfig DatabaseConfig `json:"DatabaseConfig"`
		LoggerConfig   LoggerConfig   `json:"LoggerConfig"`
	}
	DatabaseConfig struct {
		Driver   string `json:"Driver"`
		Host     string `json:"Host"`
		Port     uint16 `json:"Port"`
		Database string `json:"Database"`
		UserName string `json:"UserName"`
		Password string `json:"Password"`
	}
	LoggerConfig struct {
		Driver      string `json:"Driver"`
		Destination string `json:"Destination"`
	}
)
