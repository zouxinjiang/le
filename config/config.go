package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/zouxinjiang/le/pkgs/config"
)

var defaultFileConfig = FileConfig{
	WebConfig: WebConfig{
		Address: "",
		Port:    80,
	},
	DatabaseConfig: DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "le",
		UserName: "root",
		Password: "123456",
	},
	LogConfig: LogConfig{},
	WeiXinConfig: WeiXinConfig{
		AppId:  "wxec546224e7e2bb9c",
		Secret: "c4dd1034019479945bf590c444f6a0fc",
		Token:  "zxj123456",
	},
}

var appconf = AppConfig{
	MemoryConfig: MemoryConfig{
		AppPath:        appBase,
		ConfigFileName: "/config/config.json",
	},
	FileConfig: defaultFileConfig,
}

func Init() error {
	fconf, err := ReadConfig()
	if os.IsNotExist(err) {
		//文件不存在，则写一次文件
		_ = WriteConfig(appconf.FileConfig)
	}
	if err != nil {
		return err
	}
	fconf, _ = ReadConfig()
	appconf.FileConfig = fconf
	return nil
}

//热加载配置文件
func reloadConfig() {

}

// 写配置
func WriteConfig(fconf FileConfig) error {
	fpath := GetConfig("MemoryConfig.AppPath") + "/" + GetConfig("MemoryConfig.ConfigFileName")
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	tmp, _ := json.Marshal(fconf)
	_, _ = f.Write(tmp)
	return f.Sync()
}

// 读配置
func ReadConfig() (fconf FileConfig, err error) {
	fpath := GetConfig("MemoryConfig.AppPath") + "/" + GetConfig("MemoryConfig.ConfigFileName")
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return defaultFileConfig, err
	}
	defer f.Close()

	var data = FileConfig{}

	tmp, _ := ioutil.ReadAll(f)
	err = json.Unmarshal(tmp, &data)
	if err != nil {
		return defaultFileConfig, err
	}
	return data, nil
}

func GetConfig(key string) string {
	v, _ := GetConfigWithError(key)
	return v
}

func GetConfigWithError(key string) (string, error) {
	return config.GetConfigItem(key, appconf)
}
