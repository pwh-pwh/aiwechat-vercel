# aiwechat-vercel
使用vercel的functions，将ai功能加入微信公众号

### 介绍

无需服务器，门槛低，只需一个可以绑定到vercel的域名(无需备案)即可，基本0成本

### 快速开始

fork本项目，到vercel点击构建,环境变量填写参数

更多配置[config](conf/.env.sample)

```dotenv
# 更多配置见conf/.env.sample示例文件
GPT_TOKEN=sk-*** 你的gpt token
GPT_URL=https://xxx  代理gpt服务器(选填，默认openai官网api 例如https://api.openai.com/v1)
gptModel=gpt-3.5-turbo gpt模型(选填,默认gpt-3.5-turbo)
TOKEN=*** 微信公众号开发平台设置的token
botType=** 机器人类型 目前支持(gpt,echo,spark)例如botType=gpt
```
如何检查是否配置成功

部署后访问 vercel提供的域名/api/check 页面返回check ok即可

到域名提供商，域名增加`cname`解析到`cname-china.vercel-dns.com`

到vercel的该项目添加自定义域名(使用国内网络在访问你的域名/api/check看看能否访问)

微信公众号配置:
> 微信公众号。后台管理页面上找到`设置与开发`-`基本配置`-`服务器配置`，修改服务器地址url为`https://你的域名/api/wx` 消息加解密选择明文模式(后续添加支持加密)

### 功能支持

1. gpt回复
2. 超时回复(go协程很好用)
3. 支持连续问答(只需要在vercel创建一个redis实例，在本项目下的Storage设置连接即可，vercel会自动配置KV_URL环境变量，默认记忆对话30分钟内的内容)
4. 隐藏功能 你的域名/api/chat?msg=你的问题  (仅用于测试是否配置gpt成功,中文问题会乱码，不用管，是vercel服务器问题)
5. 检查配置：你的域名/api/check （显示当前bot的配置信息是否正确）
6. 支持图床功能，即发送图片给公众号，返回图片url

### 后续

- 支持国内大部分可以白嫖的ai 如星火(已支持，感谢大佬pr)，通义千问等
- 增加指令控制，增加管理员设置
- 增加预定义prompts
- 增加图床功能(已支持)
- 被关注自定义回复及关键词自定义回复
- 支持gemini
- 支持限制问答次数

### 杂念
项目起因:偶然看到网上有人使用vercel实现了，自己看了下文档，居然支持go了，就实现了，项目仅供学习参考
也欢迎各位大佬pr

### 问题汇总
1. 为啥要使用域名? 答: vercel提供的域名国内被墙了，微信无法访问
2. 我的域名国内可以访问，但是也不能用? 答: 这个有人反馈过，好像是微信公众平台对一些小众域名不支持

### Star History
![Star History Chart](https://api.star-history.com/svg?repos=pwh-pwh/aiwechat-vercel&type=Date)

### 项目灵感来源
[spark-wechat-vercel](https://github.com/LuhangRui/spark-wechat-vercel)
