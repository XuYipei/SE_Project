# 用户服务函数

文件位置：

*internal/net/prob.go*



## 查询更大 id 题目

`net.FindProbIdGt`

提供 id, user(opt), key(opt)

查找大于对应 id 的题目列表

显示题目的部分信息 ( 前 20 个字符的 disc )

如果给的用户 id 和 key 对得上的话就更新对应的用户在线状态

成功发送 {"success", 简略般的的题目数组}



## 查询 id 题目

`net.FindProbIdGt`

提供 id, user(opt), key(opt)

需要提供 id, email, password

往对应 email 发验证码, 建议前端搞个验证码之类的功能, 验证了之后再发邮件

注册成功发 {"success", 可忽略的用户数据}



## 查询推荐题目 

`net.FindProbRcmd`

提供 user(opt), key(opt)

推荐系统根据登录用户给出的题目

显示题目的部分信息 ( 前 20 个字符的 disc )

如果给的用户 id 和 key 对得上的话就更新对应的用户在线状态, 如果没用户暂时还没定义好行为

可能不到 maxLegth 个, bug 较多建议别显示了, 或者可以改成显示点击量排序

成功发送 {"success", 简略般的的题目数组}