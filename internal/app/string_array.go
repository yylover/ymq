package app

import "strings"

//StringArray flagset定义
type StringArray []string

//Set 实现Flag接口
func (sa *StringArray) Set(s string) error {
	*sa = append(*sa, s)
	return nil
}

func (sa *StringArray) String() string {
	return strings.Join(*sa, ",")
}
