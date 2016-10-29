calc
====

calc 是一个用于算数及逻辑表达式字符串运算的Golang库，支持自定义变量。

# 安装/更新
```bash
go get -u github.com/glace/calc
```

# 使用

## 支持的运算符以及优先级
```go
"||": optPriority1,

"&&": optPriority2,

"=":  optPriority3,
"==": optPriority3,
">=": optPriority3,
"<=": optPriority3,
">":  optPriority3,
"<":  optPriority3,

"+": optPriority4,
"-": optPriority4,

"*": optPriority5,
"/": optPriority5,
```
列表中optPriority1的优先级最低，optPriority5最高，相同优先级的运算符从左向右依次运算。"="运算符等价于"=="运算符。

## 逻辑表达式的运算
当运算数参与逻辑运算时，数值不等于0即认定为True。逻辑运算的结果True将转化为数字1，False转化为数字0。

## 普通运算表达式
```go
package main

import (
	"fmt"

	"github.com/glace/calc"
)

func main() {
	var calc calc.Calc
	expression := "1+2"

	result, err := calc.Calculate(expression)

	fmt.Printf("expression: %s\n", expression)
	fmt.Printf("result : %f", result)
	if err != nil {
		fmt.Printf(" (%v)", err.Error())
	}
	fmt.Printf("\n")
}
```
输出：
```
expression: 1+2
result : 3.000000
```

## 带变量的运算表达式
```go
var calc calc.Calc
expression := "1+2+pi"

calc.SetVariable("pi", "3.14")
result, err := calc.Calculate(expression)

fmt.Printf("expression: %s\n", expression)
fmt.Printf("result : %f", result)
if err != nil {
    fmt.Printf(" (%v)", err.Error())
}
fmt.Printf("\n")
```
输出：
```
expression: 1+2+pi
result : 6.140000
```

## 带变量但变量未赋值的运算表达式（报错）
```go
var calc calc.Calc
expression := "1+2+pi"

result, err := calc.Calculate(expression)

fmt.Printf("expression: %s\n", expression)
fmt.Printf("result : %f", result)
if err != nil {
    fmt.Printf(" (%v)", err.Error())
}
fmt.Printf("\n")
```
输出：
```
expression: 1+2+pi
result : 0.000000 (calculateRPN fail, operand expected a float number but given is 'pi' )
```

## 常用函数

### Calc.Calculate
定义：
```go
func (this *Calc) Calculate(expression string) (float64, error)
```
功能：对给定的expression表达式字符串进行计算，并获得结果及错误数据

### Calc.SetVariable
定义：
```go
func (this *Calc) SetVariable(key string, value string)
```
功能：设置变量key的值为value，用于随后的表达式计算

### Calc.CleanUpVariable
定义：
```go
func (this *Calc) CleanUpVariable()
```
功能：清除所有已设置的变量

