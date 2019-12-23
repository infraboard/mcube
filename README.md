# mcube

[![Go Report Card](https://goreportcard.com/badge/github.com/infraboard/mcube)](https://goreportcard.com/report/github.com/infraboard/mcube)
[![Release](https://img.shields.io/github/release/infraboard/mcube.svg?style=flat-square)](https://github.com/infraboard/mcube/releases)

微服务工具箱, 构建微服务中使用的工具集

+ http框架: 用于构建领域服务的路由框架, 基于httprouter进行封装
+ 异常处理: 定义API Exception
+ 日志处理: 封装zap, 用于日志处理
+ 服务注册: 服务注册组件
+ 缓存处理: 用于构建多级对象缓存
+ 加密解密: 封装cbc和ecies
+ 自定义类型: ftime方便控制时间序列化的类型, set集合
