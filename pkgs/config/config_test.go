package config

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	ca := []struct {
		A string
		B map[string]string
		C []int
		D [][]int
		E [][]string
		S struct {
			S1 string
			S2 int
			S3 struct {
				S31 int
				S32 struct {
					S321 string
				}
			}
		}
		X []interface{}
	}{
		{
			"a",
			map[string]string{
				"b1": "b1",
				"b2": "b2",
			},
			[]int{1, 2, 3},
			[][]int{{11, 22, 33}, {44, 55, 66}},
			[][]string{{"e11", "e12", "e13"}, {"e21", "e22", "e23"}},
			struct {
				S1 string
				S2 int
				S3 struct {
					S31 int
					S32 struct {
						S321 string
					}
				}
			}{
				"s1",
				42,
				struct {
					S31 int
					S32 struct {
						S321 string
					}
				}{
					123,
					struct {
						S321 string
					}{
						"222",
					},
				},
			},
			[]interface{}{
				[]int{555, 666, 777},
				[]string{"aaa", "bbb", "ccc"},
				"123",
			},
		},
	}

	//fmt.Println(getConfigItem("", ca[0]))
	//fmt.Println(getConfigItem("A", ca[0]))
	//fmt.Println(getConfigItem("B.b1", ca[0]))
	//fmt.Println(getConfigItem("C", ca[0]))
	//fmt.Println(getConfigItem("C.0", ca[0]))
	//fmt.Println(getConfigItem("C.0.11", ca[0]))
	//fmt.Println(getConfigItem("D", ca[0]))
	//fmt.Println(getConfigItem("E.1.1", ca[0]))
	//fmt.Println(getConfigItem("S.S1", ca[0]))
	//fmt.Println(getConfigItem("S.S3.S31", ca[0]))
	//fmt.Println(getConfigItem("S.S3.S32.S3216", ca[0]))
	fmt.Println(getConfigItem("X.0", ca[0]))

}
