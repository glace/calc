package calc

import (
	"fmt"
	"testing"
)

var calcTestData = []struct {
	expression                   string
	result                       float64
	isExpressionWithVarableVaild bool
	varable                      map[string]string
}{
	{"1+2", 3, true, map[string]string{}},
	{"1/2", 0.5, true, map[string]string{}},
	{"1/2+0.75", 1.25, true, map[string]string{}},
	{"1>2", 0, true, map[string]string{}},
	{"1<=2", 1, true, map[string]string{}},
	{"pv>100", 0, true, map[string]string{"pv": "50"}},
	{"pv>100&&uv<50", 0, true, map[string]string{"pv": "50", "uv": "100"}},
	{"pv>100&&uv<50", 1, true, map[string]string{"pv": "150", "uv": "40"}},
	{"pv>100&&uv<50", 0, true, map[string]string{"pv": "50", "uv": "40"}},
	{"pv>100||uv<50", 1, true, map[string]string{"pv": "50", "uv": "40"}},
	{"(pv>100)||(uv<50)", 1, true, map[string]string{"pv": "50", "uv": "40"}},
	{"(2>1)||(1+2==3)", 1, true, map[string]string{}},
	{"(2<1)||(1+2==3)", 1, true, map[string]string{}},
	{"5>4||(2<1)&&(1+2==3)", 1, true, map[string]string{}},
	{"(5>4||(2<1))&&(1+2==3)", 1, true, map[string]string{}},
	{"(5>4  ||    \t(2<1))&&  (1  +2==3)", 1, true, map[string]string{}},
	{"1+2==4", 0, true, map[string]string{}},
	{"1+2=4", 0, true, map[string]string{}},
	{"1+", 0, false, map[string]string{}},
	{"1+1)", 0, false, map[string]string{}},
	{"pv>1", 0, false, map[string]string{}},
}

func TestCalc(t *testing.T) {
	var calc *Calc

	for key, node := range calcTestData {
		fmt.Printf("Test %d :\n", key)
		fmt.Printf("    expectd %s = %f", node.expression, node.result)
		if node.isExpressionWithVarableVaild {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" (invaild)\n")
		}
		if node.varable != nil && len(node.varable) > 0 {
			fmt.Printf("    varable:\n")
			for key, value := range node.varable {
				fmt.Printf("        %s : %s\n", key, value)
			}
		}

		calc = new(Calc)

		if node.varable != nil && len(node.varable) > 0 {
			for key, value := range node.varable {
				calc.SetVariable(key, value)
			}
		}

		result, err := calc.Calculate(node.expression)

		vaild := true
		if err != nil {
			fmt.Printf("    calc result: %f(%v)\n", result, err)
			vaild = false
		} else {
			fmt.Printf("    calc result: %f\n", result)
		}

		if result == node.result && vaild == node.isExpressionWithVarableVaild {
			fmt.Printf("    statu: PASS\n")
		} else {
			fmt.Printf("    statu: FAIL\n")
			t.Fatal()
		}

		fmt.Printf("\n")
	}
}
