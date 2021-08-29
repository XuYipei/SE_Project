# 用户服务函数

文件位置：

*internal/net/user.go*



## 用户注册

`net.RegisterUser`

input:

```
json{
    name string,
    password string,
    email string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
    user json{
        id string, 
        name string,
	    email string, 
	    password string, ("")
        key string, ("")
    },
}
```

返回该用户账号的 id

注册后账号仍然无法登录, 验证状态为 false, 需要发邮件验证, 把邮件的验证码输进来

注册成功发 {"success", 有 id 字段的用户数据}



## 邮箱验证

`net.VerifyEmail`

```
json{
    id string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
    user json{
        id string,
        name string,
	    email string
	    password string, ("")
        key string, ("")
    },
}
```

需要提供 id, email, password

往对应 email 发验证码, 建议前端搞个验证码之类的功能, 验证了之后再发邮件

注册成功发 {"success", 可忽略的用户数据}



## 用户验证

input:

```
json{
    id string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
    user json{
        id string,
        name string,
	    email string
	    password string, ("")
        key string, ("")
    },
}
```

需要提供 id, password, verifyCode

查询数据库中用户的验证码是否对得上, 对得上就改变验证状态

验证成功发 {"success", 可忽略的用户数据}



## 用户登录

`net.LoginUser`

input:

```
json{
    id string,
    password string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
    user json{
        id string,
        name string,
	    email string
	    password string, ("")
        key string, 
    },
}
```

需要提供 id, password

如果已验证邮箱而且 id, password 对得上就登录成功, 生成一个登录状态

如果用户再前端进行操作就更新后端的登录状态, 超过一定时间不操作就会强制下线

暂时没有强制下线的通知, 但如果用户进行只有登录才能进行的操作 ( 推荐题目, 打分 ) 就会报错

key 值是一个用户登录状态密码, 前端发发带 key 和用户 id 的请求表示该用户对应登录信息在操作

注册成功发 {"success", 有 key 值的用户数据}



## 退出登录

`net.LogoutUser`

input:

```
json{
    id string,
    key string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
}
```


需要提供 id, password, key

删除用户对应的登录状态

退出登录成功发 {"success"}



## 删除账号

`net.RemoveUser`

input:

```
json{
    id string,
}
``` 

output:

```
json{
    status string, ("sucess" / "fail")
}
```

需要提供 id, password

删除用户( 号没了 )

删除成功发 {"success"}



