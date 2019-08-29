package config

import (
	"encoding/json"
	"fmt"
	"github.com/zouxinjiang/le/pkgs/config"
	"io/ioutil"
	"os"
	"path"
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
	reloadConfig()
	fconf, err := ReadConfig()
	if os.IsNotExist(err) {
		//文件不存在，则写一次文件
		err := WriteConfig(appconf.FileConfig)
		fmt.Println(err)
	}
	if err != nil {
		return err
	}
	fconf, _ = ReadConfig()
	appconf.FileConfig = fconf
	return nil
}

// 写配置
func WriteConfig(fconf FileConfig) error {
	fpath := GetConfig("MemoryConfig.AppPath") + "/" + GetConfig("MemoryConfig.ConfigFileName")
	_ = os.MkdirAll(path.Dir(fpath), 0777)
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	tmp, err := json.MarshalIndent(fconf, "", "\t")
	fmt.Println()
	_, _ = f.Write(tmp)
	return f.Sync()
}

// 读配置
func ReadConfig() (fconf FileConfig, err error) {
	fpath := GetConfig("MemoryConfig.AppPath") + "/" + GetConfig("MemoryConfig.ConfigFileName")
	f, err := os.OpenFile(fpath, os.O_RDWR, 0666)
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
