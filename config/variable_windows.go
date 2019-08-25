package config

import "os"

var userDir, _ = os.UserHomeDir()
var appBase = userDir + "/" + "le"
