# DropLink Backend 🖥️

## Install modules 📦

``` bash title="droplink/be"
go mod tidy
```

## Build 👷‍♂️

``` bash title="droplink/be"
go build
```

## Run development 🚀
``` bash title="/be"
go run main.go
```

## Docker

1. Build
``` bash title="/be BUILD"
docker buildx build --network=host -t droplink-be . 
```
2. Run
``` bash title="/be RUN"
docker run \
    -p $HOST_PORT:80 \
    -v $(pwd)/uploads:/app/uploads \
    -v $(pwd)/public:/app/public
    droplink-be
```