# Server

公网域名 https://accord2020.treedom.cn  
微信小程序公有服务接口  
wx.getSystemInfoSync() 返回, mpProject 服务端分配固定值默认: default  
请求头设置 // X-Custom-Dev=brand^model^system^platform^mpProject  
请求要带上小程序 X-MP-AppId  
若已登陆情况下,请求要带上身份头 X-MP-Token  

- [/leaf/v1.public.Server/Segment](#leafv1publicserversegment)
- [/leaf/v1.public.Server/Snowflake](#leafv1publicserversnowflake)

## /leaf/v1.public.Server/Segment



### Method

POST

### Request
```javascript
{
    key: "", // type:<string>
}
```

### Reply
```javascript
{
    id: "0", // type:<string(int64)>
    // Status_Success(=0) 
    // Status_Exception(=1) 
    status: "", // type:<string(enum)>
    msg: "", // type:<string>
}
```
## /leaf/v1.public.Server/Snowflake



### Method

POST

### Request
```javascript
{
    key: "", // type:<string>
}
```

### Reply
```javascript
{
    id: "0", // type:<string(int64)>
    // Status_Success(=0) 
    // Status_Exception(=1) 
    status: "", // type:<string(enum)>
    msg: "", // type:<string>
}
```
