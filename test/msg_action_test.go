package test

import (
	"testing"

	"gongsheng.cn/agent/handlers"
)

func TestProcessExternalFile_Execute(t *testing.T) {
	// 创建一个模拟的ActionInfo实例
	a := &handlers.ActionInfo{
		ctx:      nil, // 替换为你实际的上下文
		info:     &handlers.MsgInfo{},
		handlers: handlers.MessageHandler{},
	}

	// 创建ProcessExternalFile实例
	pf := &handlers.ProcessExternalFile{}

	// 调用Execute方法进行测试
	result := pf.Execute(a)

	// 验证测试结果
	if !result {
		t.Errorf("Expected Execute to return true, but got false")
	}
}
