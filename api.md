##Account账户系统
###注册: POST /account/sign_in
####args:
&nbsp;&nbsp;email 注册邮箱  
&nbsp;&nbsp;password 密码，8-25位，暂时没加密
####resp:
&nbsp;&nbsp;token 令牌
###注册: POST /account/login
####args:
&nbsp;&nbsp;email 注册邮箱  
&nbsp;&nbsp;password 密码，8-25位，暂时没加密
####resp:
&nbsp;&nbsp;token 令牌

##Error code 错误码
code:50020  msg:请重新登录  后续操作:调用login进行重新登录