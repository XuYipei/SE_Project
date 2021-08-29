# 服务端网络接口说明

ip 地址端口：http://localhost:8080

## 用户服务

### 请求报文（json）

request 1：

| 字段名     | 类型   | 意义         |
| ---------- | ------ | ------------ |
| id         | string | 用户 id      |
| password   | string | 用户密码     |
| email      | string | 用户 email   |
| name       | string | 用户名       |
| key        | string | 登录状态验证 |
| verifyCode | string | 邮箱验证     |

### 回复报文（json）

user data：

| 字段名   | 类型   | 意义       |
| -------- | ------ | ---------- |
| id       | string | 用户 id    |
| name     | string | 用户名     |
| email    | string | 用户 email |
| password | string | 目前屏蔽   |
| key      | string | 登录验证   |

replay 1：

| 字段名 | 类型   | 意义                           |
| ------ | ------ | ------------------------------ |
| status | string | 是否成功（”success“ / ”fail“） |

replay 2：

| 字段名 | 类型      | 意义     |
| ------ | --------- | -------- |
| status | string    | 是否成功 |
| user   | user data | 用户数据 |

### 网络服务

| 服务名称 | 请求类型 | url                | 请求格式  | 回复格式 |
| -------- | -------- | ------------------ | --------- | -------- |
| 用户注册 | POST     | /user/register     | request 1 | reply 1  |
| 用户登录 | POST     | /user/login        | request 1 | reply 1  |
| 退出登录 | POST     | /user/logout       | request 1 | reply 2  |
| 删除账号 | POST     | /user/remove       | request 1 | reply 2  |
| 邮箱验证 | POST     | /user/verify/email | request 1 | reply 1  |
| 用户验证 | POST     | /user/verify/user  | request 1 | reply 1  |



## 题目服务

### 请求报文格式（json）

request 1：

| 字段名    | 类型   | 意义                     |
| --------- | ------ | ------------------------ |
| id        | string | 需要查询的 id            |
| maxLength | string | 查询最大长度             |
| userId    | string | 执行操作的用户 id        |
| key       | string | 登录认证（不登录任意值） |

### 回复报文格式（json）

content data：

| 字段名 | 类型     | 意义          |
| ------ | -------- | ------------- |
| id     | string   | 对应题目的 id |
| disc   | string   | 题面描述      |
| type   | string   | 题目类型      |
| tags   | []string | 标签          |
| title  | string   | 题目名        |

problem data：

| 字段名  | 类型         | 意义                             |
| ------- | ------------ | -------------------------------- |
| id      | string       | 题目ID                           |
| content | content data | 题目内容                         |
| like    | int          | 喜欢题目的人数（暂时为0）        |
| dislike | int          | 不喜欢题目的人数（暂时为0）      |
| updTime | rfc3339      | 登录认证（目前仅登录操作时有效） |
| visit   | int          | 访问数（具体作用待定，可不要）   |

reply 1：

| 字段名  | 类型         | 意义                              |
| ------- | ------------ | --------------------------------- |
| status  | string       | "success"表示成功，"fail"表示失败 |
| problem | problem data | 题目                              |

reply 2：

| 字段名   | 类型           | 意义                              |
| -------- | -------------- | --------------------------------- |
| status   | string         | "success"表示成功，"fail"表示失败 |
| problems | []problem data | 题目数组                          |

### 网络服务

| 服务名称         | 请求类型 | url             | operation | 回复格式 |
| ---------------- | -------- | --------------- | --------- | -------- |
| 查询 id 题目     | POST     | /prob/find/id   | request 1 | reply 1  |
| 查询更大 id 题目 | POST     | /prob/find/idGt | request 1 | reply 2  |
| 查询推荐题目     | POST     | /prob/find/rcmd | request 1 | reply 2  |

**maxLength** 目前为必要字段，查询多个的时候必填



## 推荐服务

### 请求报文格式（json）

request 1：

| 字段名    | 类型    | 意义                            |
| --------- | ------- | ------------------------------- |
| userId    | string  | 用户 id                         |
| probId    | string  | 题目 id                         |
| message   | string  | 推荐评语                        |
| id        | string  | 推荐 id                         |
| maxLength | int     | 查询最大长度                    |
| key       | string  | 登录认证（不登录任意值）        |
| score     | float64 | 分数                            |
| actorId   | string  | 操作者的用户 id（不登陆任意值） |

### 回复报文格式（json）

recommend data：

| 字段名  | 类型    | 意义     |
| ------- | ------- | -------- |
| userId  | string  | 用户 id  |
| probId  | string  | 题目 id  |
| message | string  | 推荐评语 |
| id      | string  | 推荐 id  |
| score   | float64 | 分数     |

reply 1：

| 字段名    | 类型           | 意义                              |
| --------- | -------------- | --------------------------------- |
| status    | string         | "success"表示成功，"fail"表示失败 |
| recommend | recommend data | 推荐                              |

reply 2：

| 字段名     | 类型             | 意义                              |
| ---------- | ---------------- | --------------------------------- |
| status     | string           | "success"表示成功，"fail"表示失败 |
| recommends | []recommend cell | 推荐数组                          |

**网络服务**

| 服务名称             | 请求类型 | url                  |           | 请求格式 |
| -------------------- | -------- | -------------------- | --------- | -------- |
| 上传推荐             | POST     | /recmd/upload        | request 1 | reply 1  |
| 修改推荐             | POST     | /recmd/update        | request 1 | reply 1  |
| 查询题目的推荐       | POST     | /recmd/find/probId   | request 1 | reply 2  |
| 查询用户的推荐       | POST     | /recmd/find/userId   | request 1 | reply 2  |
| 查询 id 推荐         | POST     | /recmd/find/id       | request 1 | reply 1  |
| 查询大于题目id的推荐 | POST     | /recmd/find/probIdGt | request 1 | reply 2  |

**maxLength** 目前为必要字段，查询多个的时候必填

