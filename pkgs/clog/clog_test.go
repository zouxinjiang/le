package clog

import (
	"os"
	"testing"
)

func TestClog(t *testing.T) {
	lg := Clog{
		w:      os.Stdout,
		level:  Lvl_Info | Lvl_Debug | Lvl_Error | Lvl_Warning,
		format: `{"fn":"$fn","line":$ln,"data":${data}}`,
	}
	lg.dataFormat = FMT_Json
	lg.Error(map[string]string{"aaa": "ccc"})
	lg.Debug("aaa")
}

func TestCC(t *testing.T) {
	Warning("www", "tttt")
}
