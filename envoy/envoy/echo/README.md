# generate grpc code

```
protoc -I "." --go_out="server"\
    --go_opt=paths=source_relative \
    --go-grpc_out="server" --go-grpc_opt=paths=source_relative \
    echo.proto
```

# docker

```
docker build --network host -t "king011/envoy-echo:latest" .
```

```
docker push king011/envoy-echo:latest
```