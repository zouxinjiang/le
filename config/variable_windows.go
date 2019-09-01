// +build windows
package config

import (
	"os"
)

var userDir, _ = os.UserHomeDir()
var appBase = userDir + "/" + "le"

//热加载配置文件
func reloadConfig() {}
