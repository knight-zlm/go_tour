### 生成grpc go文件
` protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto `

### 生成grpc-gateway 文件
` protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. ./proto/*.proto `

### 生成swagger 文件
` go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/... `

### 通过proto 生成 swagger js
` protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. ./proto/*.proto `
