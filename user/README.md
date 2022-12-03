#### micro开发流程

1. 创建模块user
    ```shell
    docker pull micro/micro:v2.9.3
    
    #window绝对路径
    docker run --rm -v D:\go_demo\micro_shop:/www -w /www micro/micro:v2.9.3 new user
    
    #mac、linux
    docker run --rm -v $(pwd):$(pwd) -w $(pwd) micro/micro:v2.9.3 new user
    ```

2. 为user服务编写proto文件，生成proto代码
    ```shell
    # 1.先在user.proto文件增加（如果使用修改后的源码编译生成的micro，或使用修改后的micro docker生成的，则不需要，模板自动生成时会有。）
    option go_package = "./proto/user/;go_micro_service_user";
    
    # 2.然后执行在user服务根目录执行
    protoc --go_out=./ --micro_out=./ .\proto\user\user.proto
    ```

3. model层定义数据库元数据，repository层编写数据库操作，service层编写业务逻辑，handler层(相当于mvc controller)暴露api，main入口注册micro服务。
4. 编译linux平台可执行文件
   ```shell
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user-service *.go
   ```
5. 编写Dockerfile
   ```dockerfile
   FROM alpine
   ADD user-service /user-service
   ENTRYPOINT [ "/user-service" ]
   ```
6. 构建镜像
   ```shell
   docker build -t user:latest .
   ```