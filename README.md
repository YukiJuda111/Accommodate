# Accommodate
基于golang的租房网站:[项目文档](https://www.yuque.com/yuk1-note/qnreni)
## 运行环境
- go
- mysql
- redis
- consul

## 运行
```shell
# 启动consul
consul agent -dev
# 初始化mysql
mysql -uroot -p < init.sql
# 初始化redis
redis-server 
# 添加项目依赖
go mod tidy
# 启动微服务
go run service/getCaptcha/main.go
go run service/user/main.go
# 启动web服务
go run web/main.go
```

## 技术栈
- 微服务: gRPC, consul, go-micro
- web端: gin
- 数据存储: gorm+mysql, redis
