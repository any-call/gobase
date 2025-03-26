package mymail

import "testing"

func TestSend(t *testing.T) {
	content := `
<html>
<head>
    <title>Test Email</title>
</head>
<body>
    <h1>欢迎测试 HTML 邮件</h1>
    <p>这是一个用于测试的 HTML 邮件。</p>
    <p>请点击以下链接访问我们的网站：</p>
    <a href="https://www.example.com">访问我们的网站</a>
</body>
</html>
`

	err := SendByGmail("", "serverbp001@gmail.com",
		"",
		"156711203@qq.com",
		"register code11", content, nil, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send ok")
}
