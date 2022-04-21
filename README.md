<!--
 * @Description: 
 * @Author: neozhang
 * @Date: 2022-01-02 10:34:14
 * @LastEditors: neozhang
 * @LastEditTime: 2022-04-21 23:20:10
-->
# mbook  

[项目来源_性能优化+架构迭代升级，Go读书社区web开发与架构优化](https://coding.imooc.com/class/403.html)  

mbook是基于[BookStack](https://github.com/TruthHun/BookStack)进行重构和开发的。这两个项目最终来源于[MinDoc](https://github.com/lifei6671/mindoc)。
在开发的过程中，根据需求做了功能裁减，然后进行了大量架构重构和性能优化，对其中很大一部分代码的逻辑实现部分进行了改写，同时也对部分表结构进行了调整。  

先使用Go语言及Beego框架进行项目开发，快速迭代上线，然后进行包括主从和分表分库、搜索优化、页面静态化、动态缓存、下载优化、服务负载均衡等一系列架构优化，最后实现Web应用的高可用&高并发。  

## 涉及技术栈  

- Go  
- Beego框架  
- MySQL数据库和ORM  
- 前端  

## 编译运行  

1.将项目放在GOPATH/src下，使得目录结构最终如下面的样子  
```
$ tree -L 1
.
├── common
├── conf
├── controllers
├── doc
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── mbook
├── mbook.sql
├── mbook_useraction.sql
├── models
├── README.md
├── res
├── routers
├── static
├── store
├── sysinit
├── tests
├── utils
└── views
```

2.命令行到代码目录下  
```
cd $GOPATH/src/mbook
```

3.将mbook.sql和mbook_useraction.sql导入数据库  

4.运行服务  
```
bee run
```

5.在`http://localhost:8880`下登录  
![](./res/login.png)  

默认管理员用户名:`admin` , 密码:`123456`  

6.登录成功后跳转  

![](./res/demo.png)  

个人页面  

![](./res/admin.png)  

项目相关文档在[这里](/doc/)  

## 总结  

### V1.0业务快速搭建  

- 基于Beego快速搭建Web应用  
- 首页&分类&详情模块快速构建  
- 社区模块化搭建思路  
- 搜索模块快速搭建  

### V1.1~V1.5并发优化：数据库  

- 与开发相关的数据层基础优化  
- MySQL binlog与主从分离实现  
- MySQL分表分库  
- 搜索模块接入Elasticsearch  

### V2.0~V2.2并发优化：缓存层  

- 页面静态化  
- 基于Redis的动态缓存实现  
- CDN下载优化  

### V2.5并发优化：服务层  

- 代理与反向代理  
- 无状态服务与服务平行扩展  
- 负载均衡原理及其基于Nginx实践  
- 多机部署之Session同步问题  