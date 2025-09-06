# 优先级规划


## 配置区 优先级规划

+ 基础层: 9xx
+ 框架层: 8xx

+ 组件层: 6xx
+ 组合层: 5xx

+ 服务客户端: 3xx
+ 框架插件层: 2xx


## API区 优先级规划

+ API模块 -xx
+ API组件 -1xx


## trace

```sh
curl -H "traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01" http://localhost:8020/api/exapmle/v1/api.helloserviceapihandler/
```