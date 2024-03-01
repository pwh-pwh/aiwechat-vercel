# aiwechat-vercel
使用vercel的functions，将ai功能加入微信公众号

### 介绍

无需服务器，门槛低，只需一个可以绑定到vercel的域名即可，基本0成本

### 快速开始

fork本项目，到vercel点击构建,环境变量填写参数

```dotenv
GPT_TOKEN=sk-*** 你的gpt token
GPT_URL=https://xxx  代理gpt服务器(选填，默认openai官网api)
TOKEN=*** 微信公众号开发平台设置的token
TYPE=** 机器人类型 目前支持(gpt,echo)例如TYPE=gpt
```
如何检查是否配置成功
部署后访问 你的域名/api/check 页面返回check ok即可
域名增加`cname`解析到`cname-china.vercel-dns.com`

微信公众号配置:
> 微信公众号。后台管理页面上找到`设置与开发`-`基本配置`-`服务器配置`，修改服务器地址url为`https://你的域名/api/wx` 消息加解密选择明文模式(后续添加支持加密)

### 功能支持

1. gpt回复
2. 超时回复(go协程很好用)
3. 支持连续问答(todo 需要使用redis redis也可以白嫖 后续更新)
4. 隐藏功能 你的域名/api/chat?msg=你的问题 

### 后续

- 支持国内大部分可以白嫖的ai 如星火，等
- 增加记忆功能
- 增加指令控制，增加管理员设置
- 增加预定义prompts

### 杂念
项目起因:偶然看到网上有人使用vercel实现了，自己看了下文档，居然支持go了，就实现了，项目仅供学习参考
也欢迎各位大佬pr

### 项目灵感来源
[spark-wechat-vercel](https://github.com/LuhangRui/spark-wechat-vercel)
