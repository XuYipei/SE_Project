# 用户服务函数

文件位置：

*internal/net/prob.go*



## 上传推荐  

`net.UploadRecmd`

提供 userId, probId, key, message, actorId

上传推荐

需要在线的用户, 更新在线状态

成功返回 {"success", 带 id 的该推荐}



## 修改推荐

`net.UpdateRecmd`

提供 userId, probId, key, message, actorId

修改推荐

需要在线的用户, 更新在线状态

成功返回 {"success", 带 id 的该推荐}



## 查询大于题目id的推荐

`net.FindRecmdProbIdGt`

可能不展示, 没用可以不管

提供 userId, key(opt), actorId(opt)

查找对应 id 的推荐

更新在线状态

成功返回 {"success", 比对应题目 id 大的题目的推荐}



## 查询 id 推荐

`net.FindRecmdId`

提供 id, userId(opt), key(opt), actorId(opt)

查找对应 id 的推荐

更新在线状态

成功返回 {"success", 带 id 的该推荐}



## 查询用户的推荐  

`net.FindRecmdUserId`

提供 userId, key(opt), actorId(opt)

查找对应用户的推荐

更新在线状态

成功返回 {"success", 对应用户的推荐的数组}



## 查询题目的推荐

`net.FindRecmdProbId`

提供 probId,  key(opt), actorId(opt)

查找对应题目的推荐

更新在线状态

成功返回 {"success", 对应题目的推荐的数组}

