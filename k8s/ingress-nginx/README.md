# 這個例子演示了如何在 k8s 中使用 ingress-nginx 作爲 api 網關爲 http 和 grpc 提供路由功能

# server

server 是一個用於測試的 服務器與客戶端

以服務器運行
```
export ExampleAddr=127.0.0.1:8080
export ExampleMode=server
./bin/server
```

以 客戶端 運行
```
export ExampleAddr=127.0.0.1:8080
export ExampleClient=http
export ExampleClient=grpc
export ExampleMode=client
./bin/server
```