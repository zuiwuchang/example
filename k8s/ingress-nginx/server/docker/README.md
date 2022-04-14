# server
```
docker run --rm \
    --name test-server \
    -it \
    -p 9000:9000 \
    -e ExampleAddr="0.0.0.0:9000" \
    king011/k8s-ingress-nginx-example-server:0.0.1
```

# client
```
docker run --rm \
    --name test-client \
    -it \
    -e ExampleMode=client \
    -e ExampleClient=grpc \
    --network host \
    king011/k8s-ingress-nginx-example-server:0.0.1
```