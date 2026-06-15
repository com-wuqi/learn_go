package utils

import (
	"errors"
	"fmt"
	"strings"
)

type Example1 struct {
	v1 int
	v2 []string
}

func Withv1(v1 int) func(*Example1) {
	// 返回一个接收指针并修改的函数
	return func(e1 *Example1) {
		e1.v1 = v1
	}
}

func Withv2(v ...string) func(*Example1) {
	return func(e1 *Example1) {
		for _, aString := range v {
			e1.v2 = append(e1.v2, aString)
		}
	}
}

func Newexample1(opts ...func(*Example1)) *Example1 {
	e1 := &Example1{}
	for _, opt := range opts {
		opt(e1)
	}
	return e1
}

func Extra1() {
	fmt.Println("extra 1") //
	// https://chat.qwen.ai/c/79ff671a-42fe-4f6b-8b83-b0241e04d237
	// 函数式选项模式 的 构造函数
	aExample := Newexample1(Withv1(12), Withv2("hello world"))
	fmt.Println(aExample.v1, aExample.v2)
	// 链式调用实现
	aString := NewQueryBuilder()
	aString.TableName("test").Conditions("ID = 1").Limit(1)
	if v2, _error := aString.Build(); _error == nil {
		fmt.Println(v2)
	}
	// https://chat.qwen.ai/c/79ff671a-42fe-4f6b-8b83-b0241e04d237

}

type QueryBuilder struct {
	tableName  string
	conditions []string
	limitCount int
	hasLimit   bool
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		conditions: make([]string, 0),
	}
}

func (qb *QueryBuilder) TableName(table string) *QueryBuilder {
	qb.tableName = table
	return qb
}

func (qb *QueryBuilder) Conditions(conditions ...string) *QueryBuilder {
	qb.conditions = append(qb.conditions, conditions...)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitCount = limit
	qb.hasLimit = true
	return qb
}

func (qb *QueryBuilder) Build() (string, error) {
	if qb.tableName == "" {
		return "", errors.New("table name is empty")
	}
	var strBuilder strings.Builder
	strBuilder.WriteString(fmt.Sprintf("SELECT * FROM %s", qb.tableName))
	if len(qb.conditions) > 0 {
		strBuilder.WriteString(" WHERE " + strings.Join(qb.conditions, " AND "))
	}
	if qb.hasLimit {
		strBuilder.WriteString(fmt.Sprintf(" LIMIT %d", qb.limitCount))
	}
	return strBuilder.String(), nil
}
