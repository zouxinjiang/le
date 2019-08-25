package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ConfigFile = AppPath + "/conf/app_conf.json"
var appconfig AppConfig

func defaultConfig() AppConfig {
	return AppConfig{
		Inited:      false,
		ConfigDir:   AppPath + "/conf",
		TemplateDir: AppPath + "/temp",
		WebDir:      AppPath + "/web",
		BinDir:      AppPath + "/bin",
		DatabaseConfig: DatabaseConfig{
			Driver:   "",
			Host:     "",
			Port:     0,
			Database: "",
			UserName: "",
			Password: "",
		},
		LoggerConfig: LoggerConfig{
			Driver:      "os.stdout",
			Destination: "",
		},
	}
}

func LoadConfigFromFile(f string) (AppConfig, error) {
	var tmpconf AppConfig
	fp, err := os.OpenFile(f, os.O_RDONLY, 0666)
	if err != nil {
		return tmpconf, err
	}
	defer fp.Close()
	tmp, err := ioutil.ReadAll(fp)
	if err != nil {
		return tmpconf, err
	}

	err = json.Unmarshal(tmp, &tmpconf)
	if err != nil {
		return tmpconf, err
	}
	return tmpconf, nil
}

func WriteConfigToFile(conf AppConfig, f string) error {
	fp, err := os.OpenFile(f, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()
	tmp, _ := json.Marshal(conf)
	fp.Write(tmp)
	fp.Sync()
	return nil
}

//生成配置模板
func GenConfigTemplate(dst string) error {
	conf := defaultConfig()
	return WriteConfigToFile(conf, dst)
}

func Init() error {
	//加载配置文件
	conf, err := LoadConfigFromFile(ConfigFile)
	if err != nil {
		return err
	}
	appconfig = conf
	return nil
}

//=========================配置结构的操作函数=====================
func ConfigItem(key string) string {
	res, _ := ConfigItemWithError(key)
	return res
}
func ConfigItemWithError(key string) (string, error) {
	return getConfigItem(key, appconfig)
}

func GetConfigItem(key string, conf interface{}) (data string, err error) {
	return getConfigItem(key, conf)
}

func getConfigItem(key string, conf interface{}) (data string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1.(error)
			data = ""
		}
	}()
	if key == "" {
		return "", nil
	}
	karr := strings.Split(key, ".")
	k1 := karr[0]
	klen := len(karr)
	switch val := reflect.ValueOf(conf); val.Kind() {
	case reflect.String:
		if klen == 1 {
			return val.String(), nil
		}
	case reflect.Bool:
		if klen == 1 {
			return fmt.Sprintf("%v", val.Bool()), nil
		}
	case reflect.Map:
		if klen > 1 {
			val = val.MapIndex(reflect.ValueOf(k1))
			return getConfigItem(strings.Join(karr[1:], "."), val.Interface())
		} else {
			return parseMap(k1, val.Interface())
		}
	case reflect.Struct:
		if klen > 1 {
			val = val.FieldByName(k1)
			return getConfigItem(strings.Join(karr[1:], "."), val.Interface())
		} else {
			return parseStruct(k1, val.Interface())
		}
	case reflect.Chan:
		return "", errors.New("not support type of chan")
	case reflect.Slice, reflect.Array:
		if klen > 1 {
			if numReg.Match([]byte(k1)) {
				idx, _ := strconv.Atoi(k1)
				val = val.Index(idx)
			} else {
				return "", errors.New(fmt.Sprintf("slice key must number but get %v", k1))
			}
			return getConfigItem(strings.Join(karr[1:], "."), val.Interface())
		} else {
			return parseSlice(k1, val.Interface())
		}
	default:
		return fmt.Sprintf("%v", val.Interface()), nil
	}

	return
}

func parseMap(key string, m interface{}) (string, error) {
	val := reflect.ValueOf(m)
	if !val.IsValid() || val.IsNil() {
		return "", errors.New("value is not invalid")
	}
	if val.Kind() != reflect.Map {
		return "", errors.New("value is not a map")
	}
	v := val.MapIndex(reflect.ValueOf(key))
	switch v.Kind() {
	case reflect.Map, reflect.Struct, reflect.Chan:
		return "", errors.New(fmt.Sprintf("%s is not a sample type,it is type of %s", key, v.Kind()))
	case reflect.Slice, reflect.Array:
		return parseSlice("", v.Interface())
	case reflect.Invalid:
		return "", errors.New(fmt.Sprintf("%s is not found", key))
	case reflect.Interface:
		return fmt.Sprintf("%v", v.Interface()), nil
	default:
		return fmt.Sprintf("%v", v.Interface()), nil
	}
}

func parseStruct(key string, m interface{}) (string, error) {
	val := reflect.ValueOf(m)
	if !val.IsValid() {
		return "", errors.New("value is not invalid")
	}
	if val.Kind() != reflect.Struct {
		return "", errors.New(fmt.Sprintf("value is neet a struct but %s\n", val.Kind()))
	}
	v := val.FieldByName(key)
	switch v.Kind() {
	case reflect.Map, reflect.Struct, reflect.Chan, reflect.Interface:
		return "", errors.New(fmt.Sprintf("%s is not a sample type,it is type of %s", key, v.Kind()))
	case reflect.Slice, reflect.Array:
		return parseSlice("", v.Interface())
	case reflect.Invalid:
		return "", errors.New(fmt.Sprintf("key %s is not found", key))
	default:
		return fmt.Sprintf("%v", v.Interface()), nil
	}
}

var numReg = regexp.MustCompile(`^\d+$`)

func parseSlice(numKey string, m interface{}) (string, error) {
	keys := strings.Split(numKey, ".")
	k1 := keys[0]
	val := reflect.ValueOf(m)
	if !val.IsValid() || val.IsNil() {
		return "", errors.New("value is not invalid")
	}
	if val.Kind() != reflect.Slice {
		return "", errors.New("value is not a slice")
	}
	if len(keys) == 1 {
		var res = make([]string, 0)
		//数字键为空，那么返回数组成员组成的字符串
		if numKey == "" {
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i)
				if item.Kind() == reflect.Struct || item.Kind() == reflect.Map || item.Kind() == reflect.Chan {
					return "", errors.New(fmt.Sprintf("%s is not support complexe data type", val.Kind()))
				} else if item.Kind() == reflect.Slice || item.Kind() == reflect.Array {
					tmp, err := parseSlice("", item.Interface())
					if err != nil {
						return "", err
					}
					res = append(res, tmp)
				} else {
					res = append(res, fmt.Sprintf("%v", item.Interface()))
				}
			}
		} else {
			//带有数字键。那么就只返回该键对应的值
			if !numReg.Match([]byte(k1)) {
				return "", errors.New(fmt.Sprintf("slice key must number but get %v", k1))
			}
			idx, _ := strconv.Atoi(k1)
			if idx >= val.Len() {
				return "", errors.New(fmt.Sprintf("slice out of range"))
			}

			itm := val.Index(idx)
			switch itm.Kind() {
			case reflect.Slice, reflect.Array:
				return parseSlice("", itm)
			case reflect.Struct, reflect.Map, reflect.Chan, reflect.Interface:
				return "", errors.New(fmt.Sprintf("%s is not support complexe data type", val.Kind()))
			default:
				fmt.Println(itm.Kind())
				return fmt.Sprintf("%v", itm.Interface()), nil
			}
		}
		return strings.Join(res, ","), nil
	} else if numReg.Match([]byte(k1)) && len(keys) == 1 {
		idx, _ := strconv.Atoi(k1)
		item := val.Index(idx)
		return parseSlice("", item.Interface())
	} else {
		return parseSlice(strings.Join(keys[1:], "."), val.Interface())
	}

}
