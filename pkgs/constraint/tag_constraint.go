package constraint

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Valid(data interface{}) (err error) {
	// defer func() {
	// 	// 程序中出错，则不检查约束
	// 	if e := recover(); e != nil {
	// 		log.Println("验证字段:", e)
	// 		err = nil
	// 	}
	// }()

	val := reflect.ValueOf(data)
	ty := reflect.TypeOf(data)
	if val.Type().Kind() == reflect.Ptr {
		ty = reflect.TypeOf(data).Elem()
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		tag := ty.Field(i).Tag.Get("constraint")
		tag = strings.Trim(tag, " \r\n\t")
		if tag == "-" || len(tag) == 0 {
			continue
		}
		if !Constraint(tag).IsValid(fmt.Sprint(val.Field(i).Interface())) {
			return errors.New(ty.Field(i).Name)
		}
	}
	return nil
}
