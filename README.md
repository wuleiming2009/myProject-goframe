#MyProject
## Go运行版本
go1.15
## 前置准备
安装`go-imports`，并在goland->Preference->Tool->File Watcher中添加
## 依赖管理
项目使用go mod管理依赖
GOPROXY=https://goproxy.io,https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
## 运行
根目录下执行`make`命令进行编译<br>
`run.sh`可以运行项目<br>
项目根目录有docker-compose.yml，可以快速使用docker-compose完成redis和mysql的启动

## 目录结构
|_ bin&nbsp;&nbsp;&nbsp;二进制生成目录,make之后自动生成<br>
|_ src&nbsp;&nbsp;&nbsp;代码目录<br>
&nbsp;&nbsp;&nbsp;|_ cache&nbsp;&nbsp;&nbsp;&nbsp;redis缓存<br>
&nbsp;&nbsp;&nbsp;|_ common&nbsp;&nbsp;&nbsp;&nbsp;组件目录(业务或项目相关)<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ bcrypt&nbsp;&nbsp;&nbsp;&nbsp;加解密<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ finance&nbsp;&nbsp;&nbsp;&nbsp;资金工具<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ metadata&nbsp;&nbsp;&nbsp;&nbsp;中间件元数据<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ response&nbsp;&nbsp;&nbsp;&nbsp;请求响应组件<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ setting&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;启动设置组件<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ snow_flake&nbsp;&nbsp;&nbsp;&nbsp;雪花ID生成器<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ token&nbsp;&nbsp;&nbsp;&nbsp;jwt令牌<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|_ verify&nbsp;&nbsp;&nbsp;&nbsp;格式校验<br>
&nbsp;&nbsp;&nbsp;|_ conf&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;配置目录<br>
&nbsp;&nbsp;&nbsp;|_ controllers&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;mvc中的c控制器<br>
&nbsp;&nbsp;&nbsp;|_ docs&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;swagger生成的文档<br>
&nbsp;&nbsp;&nbsp;|_ lib&nbsp;&nbsp;&nbsp;&nbsp;底层库目录(与业务和项目无关)<br>
&nbsp;&nbsp;&nbsp;|_ models&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;数据层目录<br>
&nbsp;&nbsp;&nbsp;|_ proto&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;api协议层目录<br>
&nbsp;&nbsp;&nbsp;|_ routers&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;路由目录<br>
&nbsp;&nbsp;&nbsp;|_ service&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;业务逻辑目录<br>
&nbsp;&nbsp;&nbsp;|_ go.mod&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;mod依赖文件<br>
&nbsp;&nbsp;&nbsp;|_ main.go&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;程序入口文件<br>

<b>注意：lib包放的是跟业务和项目完全无关的工具包，common放的是跟项目或业务相关的方法工具等</b>
## 项目时区
全项目包括数据库所有时间采用UTC时区
## 日志
初始化： `zl := log.FromContext(ctx)`

debug信息 `zl.Debug(msg)` 格式化 `zl.Debugf`

普通信息 `zl.Info` 格式化 `zl.Infof`

错误信息 `zl.Error` 格式化 `zl.Errorf`
## 自动化api文档
`go get -u github.com/swaggo/swag/cmd/swag`  

`make swag`会自动下载`swag`工具

采用gin-swagger自动生成接口文档

地址：https://github.com/swaggo/gin-swagger

注释规范文档：https://swaggo.github.io/swaggo.io/

生成swagger文档，`make swag`

访问路由 `/swagger/index.html`


## testCurl

注册账号<br>
curl --location --request POST 'http://127.0.0.1:8000/account/sign_up' \
--data-raw '{"email":"123456@qq.com","name":"abc","password":"11111111"}'

<br>
登陆<br>
curl --location --request POST 'http://127.0.0.1:8000/account/login' \
--data-raw '{"email":"123456@qq.com","name":"abc","password":"11111111"}'