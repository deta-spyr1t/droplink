# DropLink Backend 🖥️

## Install modules 📦

``` bash title="droplink/be"
go mod tidy
```

## Build 👷‍♂️

``` bash title="droplink/be"
go build ./...
```

## Run development 🚀

### Runtime

``` bash title="/be"
go run main.go
```

### Docker

``` bash title="/be (Local storage)"
docker run -p 8080:8080 \
  -e STORAGE_DRIVER=local \
  -e LOCAL_STORAGE_PATH=/uploads \
  -v $(pwd)/uploads:/uploads \
  droplink-be
```

``` bash title="/be (AWS S3)"
docker run -p 8080:8080 \
  -e STORAGE_DRIVER=s3 \
  -e S3_BUCKET=${bucket-name} \
  -e AWS_REGION=${region} \
  -e AWS_ACCESS_KEY_ID=${access_key} \
  -e AWS_SECRET_ACCESS_KEY=${secret_key} \
  droplink-be
```