## mbook简介
mbook是基于[BookStack](https://github.com/TruthHun/BookStack)进行重构和开发的。这两个项目最终来源于[MinDoc](https://github.com/lifei6671/mindoc)。
在开发的过程中，根据需求做了功能裁减，然后进行了大量架构重构和性能优化，对其中很大一部分代码的逻辑实现部分进行了改写，同时也对部分表结构进行了调整。

## 编译运行
1.将项目放在GOPATH/src下，使得目录结构最终如下面的样子 
```
$ ls $GOPATH/src/ziyoubiancheng/mbook

LICENSE			controllers		routers			uploads
README.md		main.go			static			utils
cache			mbook.sql		store			views
common			mbook_useraction.sql	sysinit
conf			models			tests
```
2.命令行到代码目录下
```
cd $GOPATH/src/ziyoubiancheng/mbook
```

3.编译代码
```
go build
```

4.将mbook.sql和mbook_useraction.sql导入数据库
> 其中mbook_useraction.sql可以在第五章的时候再导入
5.运行服务
```
./mbook
```

6.默认管理员用户名:admin , 密码:135246


