# 项目基于Golang,MySQL和Hertz框架实现留言板功能 



## 1. 用户注册

**URL**: `/register`

**方法**: `POST`

**请求参数**:

| 参数名   | 类型  | 说明   |
| :------: | :----:| :----: |
| nickname| string|用户名|
| username |string|账号|
| password |string|密码|

**请求示例**

```json
{
    "nickname": "042",
    "username": "042",
    "password": "123123"
}
```





**响应参数：**

| 参数名  |  类型  |     说明     |
| :-----: | :----: | :----------: |
| message | string | 成功注册用户 |

```json
{
    "message": "成功注册用户"
}
```

**错误响应**

| 状态码 |      说明      |
| :----: | :------------: |
|  400   |  请求参数错误  |
|  500   | 服务器内部错误 |

## 2.用户登录

**URL:** `/login`

**方法:** `POST`

**请求参数**:

|  参数名  |  类型  | 说明 |
| :------: | :----: | :--: |
| username | string | 账号 |
| password | string | 密码 |

**请求示例:**

```json
{
    "username":"042",
    "password": "123123"
}
```



**响应参数:**

|  参数   |  类型  |   说明   |
| :-----: | :----: | :------: |
| message | string | 登陆成功 |

```json
{
    "message":"登陆成功"
}
```
**错误响应**

| 状态码 |      说明      |
| :----: | :------------: |
|  400   |  请求参数错误  |
|  500   | 服务器内部错误 |



## 3. 发表留言

**URL**: `/message`

**方法**: `POST`

**请求参数**:

|  参数名   |  类型  |       说明       |
| :-------: | :----: | :--------------: |
|  user_id  |  int   |      用户ID      |
|  content  | string |     留言内容     |
| parent_id |  int   | 父消息ID（可选） |

**请求示例**

```json
{
    "user_id": 1,
    "content": "我爱你 Ich Liebe Dich",
    "parent_id": 0
}
```

**响应参数：**

| 参数名  |  类型  |     说明     |
| :-----: | :----: | :----------: |
| message | string | 成功发表留言 |

```json
{
    "message": "成功发表留言"
}
```

**错误响应**

| 状态码 |      说明      |
| :----: | :------------: |
|  400   |  请求参数错误  |
|  500   | 服务器内部错误 |



## 4. 获取所有留言

**URL**: `/message`

**方法**: `GET`

**请求参数**: 无

**响应参数**:

|   参数名   |     类型      |    说明    |
| :--------: | :-----------: | :--------: |
|     id     |      int      |  消息的ID  |
|  user_id   |      int      |  用户的ID  |
|  content   |    string     |  消息内容  |
| created_at |    string     |  创建时间  |
| updated_at |    string     |  更新时间  |
| is_deleted |     bool      |  是否删除  |
| parent_id  | sql.NullInt64 | 父消息的ID |

**响应示例**

```json
[
    {
        "id": 1,
        "user_id": 1,
        "content": "我爱你 Ich Liebe Dich",
        "created_at": "2024-10-01T12:00:00Z",
        "updated_at": "2024-10-01T12:00:00Z",
        "is_deleted": false,
        "parent_id": {
            "Int64": 0,
            "Valid": false
        }
    }
]
```

**错误响应**

| 状态码 |      说明      |
| :----: | :------------: |
|  500   | 服务器内部错误 |



## 5. 删除留言

**URL**: `/message`

**方法**: `DELETE`

**请求参数**:

| 参数名 | 类型 |   说明   |
| :----: | :--: | :------: |
|   id   | int  | 留言的ID |

**请求示例**

```http
DELETE /message?id=1 HTTP/1.1
Host: example.com
```

**响应参数：**

| 参数名  |  类型  |     说明     |
| :-----: | :----: | :----------: |
| message | string | 成功删除留言 |

```json
{
    "message": "成功删除留言"
}
```

**错误响应**

| 状态码 |      说明      |
| :----: | :------------: |
|  400   |  请求参数错误  |
|  500   | 服务器内部错误 |

🏳 剩下的明天继续


