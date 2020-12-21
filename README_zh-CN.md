# EM-Auth

[English](./README.md) | 简体中文

## 介绍

本项目基于Etpmls-Micro开发。

本项目是一个总控制中心服务，集成用户、角色、权限的RBAC0的鉴权、自定义菜单、清除缓存、磁盘清理等功能。

![Process](docs/images/Process.jpg)

## 配置

### EM配置

参考Etpmls-Micro手册 [EM配置](https://github.com/Etpmls/Etpmls-Micro/blob/v1/README_zh-CN.md#em%E9%85%8D%E7%BD%AE)

### 网关配置

在配置文件中service-discovery=>service=>http=>tag中配置

```yaml
      tag: [
        "em.http.routers.em_AttachmentHttpService.entrypoints=web,websecure",
        "em.http.routers.em_AttachmentHttpService.rule=Host(`[YOUR_DOMAIN]`) && PathPrefix(`/api/attachment/`)",
        "em.http.routers.em_AttachmentHttpService.tls.certresolver=myresolver",
        "em.http.routers.em_AttachmentHttpService.middlewares=circuitBreaker_em-attachment@file,forwardAuth@file"
        "em.http.routers.em_AttachmentHttpService.service=em_AttachmentHttpService",

        "em.http.routers.em_AttachmentHttpService_Upload.entrypoints=web,websecure",
        "em.http.routers.em_AttachmentHttpService_Upload.rule=Host(`[YOUR_DOMAIN]`) && PathPrefix(`/storage/upload/`)",
        "em.http.routers.em_AttachmentHttpService_Upload.tls.certresolver=myresolver",
        "em.http.routers.em_AttachmentHttpService_Upload.service=em_AttachmentHttpService",

        "em.http.services.em_AttachmentHttpService.loadbalancer.passhostheader=true",
      ]
```

把[YOUR_DOMAIN]替换为你的域名

## 运行

MySQL/MariaDB
```go
go run -tags=mysql main.go
```
PostgreSQL
```go
go run -tags=postgresql main.go
```