package constraint

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Constraint string

type ConstraintType string

const (
	Cstr_Unknown   ConstraintType = ""
	Cstr_Require   ConstraintType = "require"
	Cstr_Range     ConstraintType = "range"
	Cstr_Min       ConstraintType = "min"
	Cstr_Max       ConstraintType = "max"
	Cstr_Enum      ConstraintType = "enum"
	Cstr_List      ConstraintType = "list"
	Cstr_DateTime  ConstraintType = "datetime"
	Cstr_ValueType ConstraintType = "type"
)

type ConstraintValueTypeName string

const (
	CstrValTn_Integer ConstraintValueTypeName = "integer"
	CstrValTn_Number  ConstraintValueTypeName = "number"
	CstrValTn_String  ConstraintValueTypeName = "string"
	CstrValTn_Email   ConstraintValueTypeName = "email"
	CstrValTn_Ip      ConstraintValueTypeName = "ip"
	CstrValTn_Mobile  ConstraintValueTypeName = "mobile"
)

func (self Constraint) IsValid(val string) bool {
	tp := self.getConstraintType()
	for k, v := range tp {
		switch k {
		case Cstr_Require:
			if !self.checkRequire(val) {
				return false
			}
		case Cstr_DateTime:
			if val != "" && !self.checkDateTime(v, val) {
				return false
			}
		case Cstr_Enum, Cstr_List:
			if val != "" && !self.checkEnum(v, val) {
				return false
			}
		case Cstr_Max:
			tmp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return false
			}
			if tmp != 0 && !self.checkMax(v, tmp) {
				return false
			}
		case Cstr_Min:
			tmp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return false
			}
			if tmp != 0 && !self.checkMin(v, tmp) {
				return false
			}
		case Cstr_Range:
			tmp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return false
			}
			if tmp != 0 && !self.checkRange(v, tmp) {
				return false
			}
		case Cstr_ValueType:
			if val != "" && !self.checkValueType(ConstraintValueTypeName(strings.ToLower(v)), val) {
				return false
			}
		}
	}
	return true
}

func (self Constraint) getConstraintType() map[ConstraintType]string {
	var res = map[ConstraintType]string{}
	tmp := strings.Split(string(self), ";")
	for _, item := range tmp {
		kv := strings.Split(item, ":")
		if len(kv) == 0 {
			continue
		}
		k := kv[0]
		v := ""
		if len(kv) >= 2 {
			v = kv[1]
		}
		res[ConstraintType(strings.ToLower(strings.Trim(k, " \r\n\t")))] = v
	}
	return res
}

var (
	// 数学区间表示 inf表示无穷 +inf:正无穷 -inf:服务器
	reg_range  = regexp.MustCompile(`[\(\[]{1}([-\+]?\d+\.\d+|[-\+]?\d+|-inf),([-\+]?\d+\.\d+|[-\+]?\d+|\+inf)[\)\]]{1}`)
	reg_enum   = regexp.MustCompile(`\[(\S)*\]`)
	reg_email  = regexp.MustCompile(`^([a-zA-Z0-9]+([-_]?[a-zA-Z0-9]+)*)*@([a-zA-Z0-9]*[-_]?[a-zA-Z0-9]+)+[\.][a-zA-Z]{2,3}([\.][a-zA-Z]{2})?$`)
	reg_ip     = regexp.MustCompile(`^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`)
	reg_mobile = regexp.MustCompile(`^(13[0-9]|14[5-9]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[13589])\d{8}$`)
)

func (self Constraint) checkRange(rangeExp string, val float64) bool {
	res := reg_range.FindAllStringSubmatch(rangeExp, -1)
	if len(res) == 0 {
		return true
	}
	item := res[0]
	//正确的格式
	if len(item) == 3 && len(item[0]) >= 5 {
		if item[1] == "-inf" && item[2] == "+inf" {
			return true
		}
		//判断括号类型
		switch item[0][0] {
		case '(':
			if item[1] == "-inf" {
				break
			}
			min, _ := strconv.ParseFloat(item[1], 64)
			if min >= val {
				return false
			}
		case '[':
			if item[1] == "-inf" {
				break
			}
			min, _ := strconv.ParseFloat(item[1], 64)
			if min > val {
				return false
			}
		}

		switch item[0][len(item[0])-1] {
		case ')':
			if item[2] == "+inf" {
				break
			}
			max, _ := strconv.ParseFloat(item[2], 64)
			if max <= val {
				return false
			}
		case ']':
			if item[2] == "+inf" {
				break
			}
			max, _ := strconv.ParseFloat(item[2], 64)
			if max < val {
				return false
			}
		}
	}
	return true
}

func (self Constraint) checkEnum(enum string, val interface{}) bool {
	res := reg_enum.FindAllStringSubmatch(enum, -1)
	if len(res) == 0 {
		return true
	}
	if len(res[0]) != 2 {
		return true
	}
	tmp := strings.Split(res[0][1], ",")
	for _, v := range tmp {
		if v == fmt.Sprintf("%v", val) {
			return true
		}
	}
	return false
}

func (self Constraint) checkMin(minExp string, val float64) bool {
	min, err := strconv.ParseFloat(minExp, 64)
	if err != nil {
		return true
	}
	if min > val {
		return false
	}
	return true
}

func (self Constraint) checkMax(maxExp string, val float64) bool {
	max, err := strconv.ParseFloat(maxExp, 64)
	if err != nil {
		return true
	}
	if max < val {
		return false
	}
	return true
}

func (self Constraint) checkDateTime(format string, val string) bool {
	_, err := time.Parse(format, val)
	return err == nil
}

func (self Constraint) checkValueType(typeExp ConstraintValueTypeName, val string) bool {
	switch typeExp {
	case CstrValTn_Email:
		if !self.checkEmail(val) {
			return false
		}
	case CstrValTn_Integer, CstrValTn_Number:
		if !self.checkNumber(val) {
			return false
		}
	case CstrValTn_Ip:
		if !self.checkIp(val) {
			return false
		}
	case CstrValTn_Mobile:
		if !self.checkMobile(val) {
			return false
		}
	}
	return true
}

func (self Constraint) checkRequire(val interface{}) bool {
	return !reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
}

func (self Constraint) checkNumber(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}

func (self Constraint) checkEmail(val string) bool {
	return reg_email.MatchString(val)
}

func (self Constraint) checkIp(val string) bool {
	return reg_ip.MatchString(val)
}

func (self Constraint) checkMobile(val string) bool {
	return reg_mobile.MatchString(val)
}
