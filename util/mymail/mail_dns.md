# 邮件发送与 DNS 记录的协议关系

## 文档定位

本文是 `mymail` 的公共基础文档，用于建立一套稳定的邮件发送模型：一封邮件如何从发送请求进入 SMTP 系统，
如何经过 ZoneMTA 投递到收件方，以及 SMTP 字段和 DNS 记录分别承担什么职责。

本文不把邮件相关概念简单罗列，而是按照以下主线组织：

```text
业务邮件
    -> SMTP Envelope 与邮件 Header
    -> SMTP Relay / ZoneMTA 投递
    -> 收件方身份检查
    -> SPF / DKIM / DMARC / PTR 结果
    -> 收件、隔离或拒收
```

理解这条主线后，再看具体 DNS 记录和 Go 代码，才能知道每个字段为什么存在、由谁提供、在哪个阶段生效。

## 一、先建立三个边界

邮件系统中有三个容易混淆的边界：

### 1. 客户端接入与公网投递

```text
smtp.example.com:587
```

是客户端连接邮件平台的入口；

```text
127.0.0.1:2525
```

是 MailGo 与本机 ZoneMTA 之间的内部提交入口；

```text
公网 IP -> 收件方 MX:25
```

才是 ZoneMTA 向 QQ、Gmail、Outlook 等收件系统进行的实际公网投递。

### 2. 邮件内容身份与 SMTP 投递身份

邮件 Header 中的 `From` 是收件人看到的发件人；SMTP Envelope 中的 `MAIL FROM` 是投递和退信使用的身份。两者可以相同，也可以不同，因此不能把它们当作同一个字段处理。

### 3. 发送授权与收信路由

SPF、DKIM、DMARC 主要用于证明“谁可以代表这个域名发信”；MX 主要用于说明“这个域名的邮件应该由谁接收”。发信授权和收信路由是两件不同的事情。

## 二、完整发送链路

以 MailGo 和 ZoneMTA 为例：

```text
SMTP 客户端
    |
    | smtp.example.com:587
    v
MailGo SMTP Relay
    |
    | 127.0.0.1:2525
    v
ZoneMTA
    |
    | 使用服务器公网 IP 连接收件方 MX 的 25 端口
    v
QQ / Gmail / Outlook
    |
    | 检查 SPF、DKIM、DMARC、PTR、IP 信誉和邮件内容
    v
收件箱、垃圾箱或拒收
```

`127.0.0.1:2525` 是发信服务器内部的 SMTP feeder 端口。收件方看不到这个地址，也不会因为邮件从 2525 端口提交，就认为邮件来自 2525。

收件方实际看到的是：

- ZoneMTA 对外连接时使用的公网 IP。
- SMTP 会话中的 `EHLO` 主机名。
- 邮件的 `MAIL FROM`，也就是 Envelope From。
- 邮件 Header 中的 `From`、`To`、`Subject` 和 DKIM 签名。
- 与这些身份对应的 DNS 记录和服务器信誉。

## 三、SMTP Envelope 与邮件 Header

SMTP 投递包含两类信息。

### 1. SMTP Envelope

```text
MAIL FROM:<bounce@example.com>
RCPT TO:<user@example.net>
```

Envelope 是 SMTP 服务器在投递过程中使用的信息，主要用于路由和退信。

- `MAIL FROM`：失败投递时的退信地址，也称 Envelope From。
- `RCPT TO`：本次 SMTP 投递的实际收件人。
- SPF 主要检查 `MAIL FROM` 所属域名是否授权当前发信 IP。

### 2. 邮件 Header

```text
From: "Example" <notice@example.com>
To: user@example.net
Reply-To: support@example.com
Subject: Welcome
```

Header 是收件人看到的邮件信息。

- `From`：用户界面中显示的发件人。
- `To`：显示给收件人的目标地址。
- `Reply-To`：用户点击回复时使用的地址。
- `Subject`：邮件标题。
- `DKIM-Signature`：发信服务器对邮件内容和指定 Header 生成的签名。

`From` 和 `MAIL FROM` 可以相同，也可以不同：

```text
From: notice@example.com
MAIL FROM:<bounce@example.com>
```

生产系统需要明确两者的关系，否则 SPF、DKIM 与 DMARC 可能无法对齐。

## 四、DNS 记录的协议角色

| 记录 | 作用 | 主要解决的问题 |
| --- | --- | --- |
| A / AAAA | 将主机名指向 IP 地址 | 客户端如何找到 SMTP 服务或发信主机 |
| MX | 指定域名的收信服务器 | 别人如何向该域名投递邮件 |
| PTR | 将公网 IP 反向解析为主机名 | 发信服务器身份和主机名信誉 |
| SPF | 授权哪些 IP 代表域名发信 | `MAIL FROM` 是否允许当前 IP |
| DKIM | 对邮件内容进行签名 | 邮件是否来自持有私钥的发信系统、内容是否被篡改 |
| DMARC | 规定 SPF、DKIM 与 Header From 的对齐和处理策略 | 收件方如何处理疑似伪造邮件 |

### A / AAAA

例如：

```text
smtp.example.com  A  203.0.113.10
```

它只表示客户端可以通过该主机名找到服务器。它不代表邮件已经通过 SPF、DKIM 或 DMARC。

对 SMTP 主机使用 Cloudflare 时，应使用 DNS Only，不能使用代理转发 SMTP 流量。

### MX

MX 用于接收邮件，不是发信授权记录：

```text
example.com  MX  10 inbound.example.com
```

不要为了发送邮件而随意修改根域名 MX。错误的 MX 记录可能导致现有收件服务失效。

如果系统暂时没有收信或退信处理服务，不应贸然把根域名 MX 指向 MailGo。

### PTR

建议保持正向和反向解析一致：

```text
203.0.113.10 -> smtp.example.com
smtp.example.com -> 203.0.113.10
```

PTR 通常需要在云服务器或 IP 服务商处设置，不能通过普通域名 DNS 面板设置。

## 五、SPF：授权发信来源

SPF 是发布在域名下的 TXT 记录，用于声明哪些服务器可以代表该域名发信。

示例：

```text
example.com  TXT  "v=spf1 ip4:203.0.113.10 -all"
```

含义：允许 `203.0.113.10` 代表 `example.com` 发信，其他来源全部失败。

注意事项：

1. 一个域名只能有一条 `v=spf1` 记录。
2. 如果同时使用其他邮件平台，需要合并 `include:` 或 `ip4:`，不能创建第二条 SPF。
3. SPF 检查的主要对象是 `MAIL FROM` / Return-Path 域名，不是邮件 Header 中显示的 `From`。
4. 初期可以使用 `~all` 观察，确认来源完整后再使用 `-all`。

## 六、DKIM：证明邮件来源并保护内容

DKIM 使用非对称密钥：

- 私钥保存在 ZoneMTA 或发信服务中，用于签名。
- 公钥发布在 DNS，用于收件方验证签名。

记录格式示例：

```text
mailgo._domainkey.example.com  TXT  "v=DKIM1; k=rsa; p=公钥内容"
```

其中 `mailgo` 是 selector 示例，必须使用 ZoneMTA 实际配置的 selector。

DKIM 公钥不能自行编造，必须与 ZoneMTA 使用的私钥成对。私钥不得写入 DNS、代码仓库、公共文档或客户端配置。

## 七、DMARC：约束身份对齐和处理策略

DMARC 通过 `_dmarc` TXT 记录发布策略：

```text
_dmarc.example.com  TXT  "v=DMARC1; p=none; rua=mailto:dmarc@example.com"
```

常用策略：

- `p=none`：只收集报告，不拦截。
- `p=quarantine`：验证失败的邮件进入垃圾箱或隔离区。
- `p=reject`：验证失败的邮件建议直接拒收。

建议上线初期使用 `p=none`，确认 SPF、DKIM 和域名对齐正常后，再逐步提高策略强度。

DMARC 主要关注 Header From 域名，并要求它和 SPF 或 DKIM 的认证域名满足对齐关系。因此下面三者不能只看其中一个：

```text
Header From 域名
MAIL FROM 域名
DKIM d= 域名
```

## 八、配置值的来源和责任边界

| 配置项 | 值的来源 |
| --- | --- |
| SMTP 服务 A 记录 | 系统管理员和服务器公网 IP |
| SPF | 根据实际公网发信 IP 和第三方发信平台整理 |
| DKIM selector | ZoneMTA 或系统配置 |
| DKIM 公钥 | 由 ZoneMTA 私钥生成 |
| DKIM 私钥 | 仅保存在发信服务器 |
| DMARC 策略 | 平台管理员根据业务风险定义 |
| PTR | 云服务器或 IP 服务商 |
| MX | 实际收信服务提供方 |
| Return-Path 域名 | 邮件系统设计和退信处理逻辑 |

## 九、MailGo 示例

以下仅为示例，真实 DKIM selector 和公钥必须从 ZoneMTA 配置中取得：

```text
smtp.mailgo.fyi       A    43.227.112.76
mailgo.fyi            TXT  "v=spf1 ip4:43.227.112.76 -all"
mailgo._domainkey     TXT  "v=DKIM1; k=rsa; p=实际公钥"
_dmarc                TXT  "v=DMARC1; p=none"
```

如果还没有收信或退信处理服务，不要为了发信而随意增加根域名 MX 记录。

## 十、验证方法

```bash
dig +short A smtp.example.com
dig +short MX example.com
dig +short TXT example.com
dig +short TXT mailgo._domainkey.example.com
dig +short TXT _dmarc.example.com
dig -x 203.0.113.10 +short
```

检查邮件是否真正通过认证，应查看收件方的“查看原始邮件”或邮件 Header，重点关注：

```text
Authentication-Results:
spf=pass
dkim=pass
dmarc=pass
```

邮件能够进入收件箱，只能说明收件方接受了这次投递，不代表 SPF、DKIM、DMARC 都已经通过。

## 十一、实现时的最小原则

为了让协议边界在代码中保持清晰，邮件工具和业务服务应遵循以下原则：

1. SMTP Envelope、邮件 Header、邮件正文分别建模，不使用一个字符串字段代替全部信息。
2. `MAIL FROM` 用于退信和 SPF 判断，`From` 用于展示和 DMARC 对齐，代码中不能混用。
3. ZoneMTA 的内部 feeder 端口只允许本机或受控内网访问，不直接暴露到公网。
4. SMTP 客户端入口使用 587 和 STARTTLS；公网投递由 ZoneMTA 负责。
5. SPF 记录只能有一条，DKIM selector 和公钥必须与实际私钥匹配，DMARC 初期使用 `p=none` 观察。
6. DNS 验证成功不等于邮件一定进入收件箱，收件方还会检查 PTR、IP 信誉、内容和投递行为。
7. 测试代码不能默认真实发信，真实投递必须通过显式环境变量或独立集成测试开启。

## 十二、排查顺序

当邮件发送失败或进入垃圾箱时，按下面顺序排查：

1. SMTP 客户端是否能连接 587。
2. MailGo 是否能连接本机 ZoneMTA feeder。
3. ZoneMTA 是否成功解析收件方 MX。
4. ZoneMTA 是否能连接收件方的 25 端口。
5. 公网 IP 的 PTR 是否正确。
6. SPF 是否授权实际的公网发信 IP。
7. DKIM 是否使用正确 selector 和对应公钥。
8. DMARC 的 Header From、SPF 域名和 DKIM `d=` 域名是否对齐。
9. 最后检查收件方的垃圾邮件策略、IP 信誉和邮件内容。
