## server
本层主要是提供服务能力如http,grpc.并把具体的实现类注入，http,grpc.进行路由绑定。
```shell
# 开发的路径
先定义proto接口定义文件->protoc 生成对应的golang代码【或者使用make api 】-> service目录下 实现具体的接口业务逻辑-> 本层 启动服务能力，并且把具体的实现与服务做绑定，也就是注册路由
#
```
