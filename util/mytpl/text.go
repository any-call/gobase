package mytpl

import (
	"bytes"
	"text/template"
)

// 定义 .tpl 文件: 创建模板文件，例如 example.tpl，其中包含占位符变量，通常用双花括号 {{ .VariableName }} 表示。

type TextTPL struct {
	*template.Template
}

func NewTextTPL(name string, TxtTpl string) (*TextTPL, error) {
	tpl, err := template.New(name).Parse(TxtTpl)
	if err != nil {
		return nil, err
	}

	return &TextTPL{tpl}, nil
}

func (self *TextTPL) RenderTemplate(data any) (string, error) {
	// 用于存储渲染后的模板内容
	var buf bytes.Buffer
	// 渲染模板，将结果写入 buffer
	if err := self.Execute(&buf, data); err != nil {
		return "", err
	}

	// 返回渲染后的字符串
	return buf.String(), nil
}
