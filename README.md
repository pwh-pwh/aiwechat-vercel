# aiwechat-vercel
使用[vercel](https://vercel.com/dashboard)的functions，将ai功能加入微信公众号
![image](https://private-user-images.githubusercontent.com/113298285/340937333-b9d80b58-3a8f-4f60-994b-b4b2c12e73ac.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3MTg3ODQ3NjcsIm5iZiI6MTcxODc4NDQ2NywicGF0aCI6Ii8xMTMyOTgyODUvMzQwOTM3MzMzLWI5ZDgwYjU4LTNhOGYtNGY2MC05OTRiLWI0YjJjMTJlNzNhYy5wbmc_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjQwNjE5JTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI0MDYxOVQwODA3NDdaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT1lOTJmZWRkNzNhMzZjYzAzMzcyMWZkMGYzZGY0NDZjMWQ4YTFiMDdjMTU3MTYyYTNhOTAzOWYzNWY5MTY0ZjA5JlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCZhY3Rvcl9pZD0wJmtleV9pZD0wJnJlcG9faWQ9MCJ9.ZfayDmDOgNi-iBZXTWaLt_H4kZCXOYR7SFpqZJwwA5s)

### 介绍

无需服务器，门槛低，只需一个可以绑定到vercel的域名(无需备案)即可，基本0成本

### 快速开始

- 提前到Vercel创建Redis数据库
  - Dashboard--Storage--Create--KV Redis
![image](https://github.com/pwh-pwh/aiwechat-vercel/assets/24500520/523f9aec-faaf-415a-8cb9-520041a06d58)

- Fork本Github项目，到Vercel点击构建,环境变量填写参数
- 在vercel该项目详情页面的Storage选择连接前面创建的redis数据库
  - 数据库链接成功后，Vercel会自动配置KV_URL环境变量

#### 数据库配置详情

图片步骤:
> ![config](http://mmbiz.qpic.cn/mmbiz_jpg/6q5SCtonIfFYZpvZdOUbibQBicXkllyO3K6XOp2zUv6PE3e1tqpfYA7wSYRWByZX9Wibibq9PDr6ML4iaiacTWNAaI9Q/0)

更多配置[config](conf/.env.sample)

```dotenv
GPT_TOKEN=sk-*** 你的gpt token
GPT_URL=https://xxx/v1  代理gpt服务器(选填，默认openai官网api 例如https://api.openai.com/v1)
gptModel=gpt-3.5-turbo gpt模型(选填,默认gpt-3.5-turbo)
WX_TOKEN=*** 微信公众号开发平台设置的token
botType=** 机器人类型 目前支持(gpt,echo,spark,qwen,gemini)例如botType=gpt
```
如何检查是否配置成功

部署后访问 vercel提供的域名/api/check 页面返回check ok即可

到域名提供商，域名增加`cname`解析到`cname-china.vercel-dns.com`

到vercel的该项目添加自定义域名(使用国内网络在访问你的域名/api/check看看能否访问)

微信公众号配置:
> 微信公众号。[微信公众平台](https://mp.weixin.qq.com/)后台管理页面上找到`设置与开发`-`基本配置`-`服务器配置`，修改服务器地址url为`https://你的域名/api/wx` 消息加解密选择明文模式(后续添加支持加密)

录制了一期简单的视频教程供参考[b站](https://b23.tv/BNWDKu1)

也有大佬写了自己在cloudflare部署的教程[discussions](https://github.com/pwh-pwh/aiwechat-vercel/discussions/22)

### 功能支持

1. 支持接入gpt,星火,通义千问,gemini
2. 超时回复(go协程很好用)
3. 支持连续问答(只需要在vercel创建一个redis实例，在本项目下的Storage设置连接即可，vercel会自动配置KV_URL环境变量，默认记忆对话30分钟内的内容)
4. 隐藏功能 你的域名/api/chat?msg=你的问题  (仅用于测试是否配置gpt成功,也可用作于简单的接口api,中文乱码问题已修复)
5. 检查配置：你的域名/api/check （显示当前bot的配置信息是否正确）
6. 支持图床功能，即发送图片给公众号，返回图片url
7. 被关注自定义回复
8. 支持设置system prompt
9. 支持指令

### 指令支持
1. /help：查看帮助
2. /gpt：切换与GPT对话
3. /spark：切换与星火对话
4. /qwen：切换与通义千问对话
5. /gemini：切换与gemini对话
6. /prompt: 你的prompt: 设置system prompt
7. /getpt: 获取当前设置prompt
8. /cpt: 清除当前设置prompt
9. /setmodel model_name:设置当前bot使用的模型
10. /setmodel:重置当前bot的模型为默认值
11. /getmodel:获取当前bot自定义的模型名
12. /clear:清除对话列表

有其它想要支持的指令欢迎提issue或者pr (例如查看天气啥的)

### 后续计划添加指令
1. /fy: 翻译文本
2. /wec: 查看天气
3. todolist管理： /ta: 添加待办事项 /td: 删除待办事项 /tl: 查看待办事项列表

### 后续

- 支持国内大部分可以白嫖的ai 如星火(已支持，感谢大佬pr)，通义千问(已支持，感谢大佬pr)等(有想要添加的可以提个issue)
- 增加指令控制(已支持)，增加管理员设置
- 关键词自定义回复
- 支持限制问答次数
- 支持企业微信群机器人
- todolist功能，用户可以在机器人管理待办事件
- 查看股票和币价

### 杂念
项目起因:偶然看到网上有人使用vercel实现了,但是功能比较单一，看了一下文档，支持go所以就想自己开发下，支持接入多一点ai和自定义功能，项目仅供学习参考
也欢迎各位大佬pr,来个免费的<mark>star<mark>

### 问题汇总
1. 为啥要使用域名? 答: vercel提供的域名国内被墙了，微信无法访问
2. 为啥有时候可以回复，有时候没有回复？答: 微信公众号限制答复500多字，超过回复会失败，可以增加限制字数的提示词解决。还有一个原因是答复太久，接口超时了免费版vercel的functions限制接口10s
3. 域名需要备案吗?答:不需要，另外也可以在cloudflare托管域名(白嫖一些2级域名，托管上去，可以达到0成本)
4. 我的是订阅号支持吗?答:无论是公众号还是订阅号,自动回复都是一个机制，所以都支持
5. 发送信息返回错误error, status code: 403, message: invalid character '<' looking for beginning of value怎么回事?答:检查GPT_URL是不是漏了/v1或者cf开了盾，墙之类的
6. 支持接入deepseek吗?答:支持，不过有一点要注意deepseek支持的模型为deepseek-coder，deepseek-chat要正常使用，需要改gptModel为这两个模型之一
7. 修改环境变量后，还是不成功?答:在修改环境变量后要重新部署下配置才后生效，因为vercel原来的实例没有被销毁读取的还是未修改的环境变量。建议每次修改环境变量后手动重新部署一下
8. 微信字数限制如何解决?答:已经有大佬提pr了，可以通过设置最大token解决，设置环境变量maxOutput即可，一般设置到500，回答没有完整可以和ai说继续即可，pr详情[pr](https://github.com/pwh-pwh/aiwechat-vercel/pull/36)

更多功能探讨[discussions](https://github.com/pwh-pwh/aiwechat-vercel/discussions)

### Star History
![Star History Chart](https://api.star-history.com/svg?repos=pwh-pwh/aiwechat-vercel&type=Date)

### 项目灵感来源
[spark-wechat-vercel](https://github.com/LuhangRui/spark-wechat-vercel)
