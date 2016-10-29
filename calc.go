package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Calc struct {
	variables map[string]string
}

// return the priority of given operator
// error will be not nil if given operator is not support
func (this *Calc) GetOptPriority(operator string) (int, error) {
	const (
		optPriority0 = iota
		optPriority1
		optPriority2
		optPriority3
		optPriority4
		optPriority5
	)

	var optPriority map[string]int = map[string]int{
		"||": optPriority1,

		"&&": optPriority2,

		"=":  optPriority3, // the same as ==
		"==": optPriority3,
		">=": optPriority3,
		"<=": optPriority3,
		">":  optPriority3,
		"<":  optPriority3,

		"+": optPriority4,
		"-": optPriority4,

		"*": optPriority5,
		"/": optPriority5,
	}

	if _, ok := optPriority[operator]; ok {
		return optPriority[operator], nil
	} else {
		return optPriority0, errors.New(fmt.Sprintf("operator '%s' not support", operator))
	}
}

// calculate the simple expression "a operator b"
// eg. normalCalculate(1,2,"+") returns 3, nil
func (this *Calc) normalCalculate(a, b float64, operator string) (float64, error) {
	var (
		result        float64          = 0
		err           error            = nil
		boolToFloat64 map[bool]float64 = map[bool]float64{
			true:  1,
			false: 0,
		}
	)
	switch operator {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	case ">":
		result = boolToFloat64[a > b]
	case "<":
		result = boolToFloat64[a < b]
	case "=":
		result = boolToFloat64[a == b]
	case "==":
		result = boolToFloat64[a == b]
	case ">=":
		result = boolToFloat64[a >= b]
	case "<=":
		result = boolToFloat64[a <= b]
	case "&&":
		result = boolToFloat64[a != 0 && b != 0]
	case "||":
		result = boolToFloat64[a != 0 || b != 0]
	default:
		err = errors.New(fmt.Sprintf("operator '%s' not support", operator))
	}
	return result, err
}

// set variable for expression
func (this *Calc) SetVariable(key string, value string) {
	if this.variables == nil {
		this.variables = make(map[string]string)
	}
	this.variables[key] = value
}

// remove all variables setted for expression
func (this *Calc) CleanUpVariable() {
	if this.variables == nil {
		this.variables = make(map[string]string)
	}
	for key, _ := range this.variables {
		delete(this.variables, key)
	}
}

// calculate the given expression string
// return the result
// eg. Calculate("1+2") returns 3, nil
func (this *Calc) Calculate(expression string) (float64, error) {
	var spiltedStr []string
	var err error
	var result float64

	spiltedStr, err = this.generateRPN(expression)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("generateRPN fail, %s", err.Error()))
	}

	spiltedStr = this.replaceVar(spiltedStr)

	result, err = this.calculateRPN(spiltedStr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("calculateRPN fail, %s", err.Error()))
	}

	return result, nil
}

func (this *Calc) replaceVar(spiltedStr []string) []string {
	var outSpiltedStr []string
	for _, value := range spiltedStr {
		if _, ok := this.variables[value]; ok {
			outSpiltedStr = append(outSpiltedStr, this.variables[value])
		} else {
			outSpiltedStr = append(outSpiltedStr, value)
		}
	}

	return outSpiltedStr
}

// calculate the Reverse Polish Notation
// return the result of given Reverse Polish Notation which is generate by generateRPN() function
func (this *Calc) calculateRPN(rpn []string) (float64, error) {
	var stack LinkStack
	stack.Init()
	for i := 0; i < len(rpn); i++ {
		if caleIsOperand(rpn[i]) {
			if f, err := strconv.ParseFloat(rpn[i], 64); err != nil {
				return 0, errors.New(fmt.Sprintf("operand expected a float number but given is '%s' ", rpn[i]))
			} else {
				stack.Push(f)
			}
		} else {
			if stack.Count < 2 {
				return 0, errors.New("operand not enough")
			}

			p1 := stack.Pop().(float64)
			p2 := stack.Pop().(float64)

			p3, err := this.normalCalculate(p2, p1, rpn[i])
			if err != nil {
				return 0, errors.New(fmt.Sprintf("normalCalculate fail, %s", err.Error()))
			}

			stack.Push(p3)
		}
	}

	if stack.Count > 1 {
		return 0, errors.New("operator not enough")
	} else if stack.Count == 0 {
		return 0, errors.New("operand not enough")
	}

	result := stack.Pop().(float64)
	return result, nil
}

func (this *Calc) generateRPN(exp string) ([]string, error) {
	var stack LinkStack
	stack.Init()

	var spiltedStr []string = this.splitString(exp)
	var rpn []string

	for i := 0; i < len(spiltedStr); i++ { // 遍历每一个元素
		curSpilt := spiltedStr[i] //当前元素

		if !caleIsOperand(curSpilt) {
			// 如果不是操作数

			// 四种情况入栈
			// 1 当前元素为左括号直接入栈
			// 2 栈内为空直接入栈
			// 3 栈顶为左括号，且当前元素不是右括号，直接入栈
			// 4 当前元素不为右括号时，在比较栈顶元素与当前元素的运算符优先级，如果当前元素大，直接入栈。
			if curSpilt == "(" ||
				stack.LookTop() == nil ||
				(stack.LookTop().(string) == "(" && curSpilt != ")") ||
				(this.compareOperator(curSpilt, stack.LookTop().(string)) == 1 && curSpilt != ")") {
				stack.Push(curSpilt)
			} else {
				if curSpilt == ")" { //当前元素为右括号时，提取操作符，直到碰见左括号
					for {
						if stack.Count <= 0 {
							return rpn, errors.New("operator '(' expected but not found")
						}
						if pop := stack.Pop().(string); pop == "(" {
							break
						} else {
							rpn = append(rpn, pop)
						}
					}
				} else { //当前元素为操作符时，不断地与栈顶元素比较直到遇到比自己小的（或者栈空了），然后入栈。
					for {
						pop := stack.LookTop()
						if pop != nil && pop != "(" && this.compareOperator(curSpilt, pop.(string)) != 1 {
							rpn = append(rpn, stack.Pop().(string))
						} else {
							stack.Push(curSpilt)
							break
						}
					}
				}
			}
		} else {
			// 如果是操作数，直接添加到后缀表达式中
			rpn = append(rpn, curSpilt)
		}
	}

	//将栈内剩余的操作符全部弹出。
	for {
		if pop := stack.Pop(); pop != nil {
			rpn = append(rpn, pop.(string))
		} else {
			break
		}
	}
	return rpn, nil
}

// compare the priority of two given operators
// if return 1, o1Priority > o2Priority
// if return 0, o1Priority = o2Priority
// if return -1, o1Priority < o2Priority
func (this *Calc) compareOperator(operator1, operator2 string) int {
	o1Priority, _ := this.GetOptPriority(operator1)
	o2Priority, _ := this.GetOptPriority(operator2)

	if o1Priority > o2Priority {
		return 1
	} else if o1Priority == o2Priority {
		return 0
	} else {
		return -1
	}
}

// split the given expression string to string array consist of Operand and Operator
// eg. 5>4||(2<1)&&(1+2==3) split to [5 > 4 || ( 2 < 1 ) && ( 1 + 2 == 3 )]
func (this *Calc) splitString(expression string) []string {
	// remove all space
	expression = strings.Replace(expression, " ", "", -1)
	// remove all tab
	expression = strings.Replace(expression, "\t", "", -1)

	var splitArr []string
	byteExp := []byte(expression)

	var operand string
	var operator string
	var curByte byte

	for i := 0; i < len(byteExp); i++ {
		curByte = byteExp[i]

		if caleIsOperandByte(curByte) {
			operand += string(curByte)
		} else {
			completeOpt := false
			if curByte == '(' || curByte == ')' {
				if operand != "" {
					splitArr = append(splitArr, operand)
					operand = ""
				}
				if operator != "" {
					splitArr = append(splitArr, operator)
					operator = ""
				}
				completeOpt = true
			}

			operator += string(curByte)

			if i+1 >= len(byteExp) {
				completeOpt = true
			}
			if i+1 < len(byteExp) && (caleIsOperandByte(byteExp[i+1]) || calcIsSpace(byteExp[i+1])) {
				completeOpt = true
			}

			if completeOpt {
				if operand != "" {
					splitArr = append(splitArr, operand)
					operand = ""
				}
				if operator != "" {
					splitArr = append(splitArr, operator)
					operator = ""
				}
			}
		}
	}
	if operand != "" {
		splitArr = append(splitArr, operand)
		operand = ""
	}
	if operator != "" {
		splitArr = append(splitArr, operator)
		operator = ""
	}
	return splitArr
}

func caleIsOperandByte(ch byte) bool {
	if calcIsAlpha(ch) || calcIsNumber(ch) {
		return true
	}
	return false
}

// 判断是否是合法的操作数（数字或变量）
// return a bool if given operand string is a valid operand(number or variable)
func caleIsOperand(operand string) bool {
	operand = strings.TrimSpace(operand)

	// 判断是否是数字
	// judge if given operand string is number
	if _, err := strconv.ParseFloat(operand, 64); err == nil {
		return true
	}

	// 判断是否是潜在的变量
	// 变量名只支持大小写字母和数字
	// judge if given operand string is variable
	// a vaild variable name here only contains [A..Za..z0..9]
	for _, v := range []byte(operand) {
		if !(calcIsAlpha(v) || (v >= '0' && v <= '9')) {
			return false
		}
	}
	return true
}

func calcIsNumber(o1 byte) bool {
	if o1 >= '0' && o1 <= '9' || o1 == '.' {
		return true
	}
	return false
}

func calcIsAlpha(o1 byte) bool {
	if o1 >= 'a' && o1 <= 'z' || o1 >= 'A' && o1 <= 'Z' {
		return true
	}
	return false
}

func calcIsSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

////////////////////////////

// stack used for calc
type StackNode struct {
	Data interface{}
	next *StackNode
}

type LinkStack struct {
	top   *StackNode
	Count int
}

func (this *LinkStack) Init() {
	this.top = nil
	this.Count = 0
}

func (this *LinkStack) Push(data interface{}) {
	var node *StackNode = new(StackNode)
	node.Data = data
	node.next = this.top
	this.top = node
	this.Count++
}

func (this *LinkStack) Pop() interface{} {
	if this.top == nil {
		return nil
	}
	returnData := this.top.Data
	this.top = this.top.next
	this.Count--
	return returnData
}

//Look up the top element in the stack, but not pop.
func (this *LinkStack) LookTop() interface{} {
	if this.top == nil {
		return nil
	}
	return this.top.Data
}
