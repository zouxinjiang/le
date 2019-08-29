package config

var appBase = "/usr/local/le"

//热加载配置文件
func reloadConfig() {
	var sign = make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGUSR1)
	go func() {
		for {
			<-sign
			ReadConfig()
		}
	}()
}
