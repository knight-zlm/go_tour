### 生成grpc go文件
` protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto `

### 生成grpc-gateway 文件
` protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. ./proto/*.proto `

### 生成swagger 文件
` go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/... `

### 通过proto 生成 swagger js
` protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. ./proto/*.proto `

### etcd docker 安装
`
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
   docker rmi gcr.io/etcd-development/etcd:v3.3.25 || true && \
   docker run \
   -p 2379:2379 \
   -p 2380:2380 \
   --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
   --name etcd-gcr-v3.3.25 \
   gcr.io/etcd-development/etcd:v3.3.25 \
   /usr/local/bin/etcd \
   --name s1 \
   --data-dir /etcd-data \
   --listen-client-urls http://0.0.0.0:2379 \
   --advertise-client-urls http://0.0.0.0:2379 \
   --listen-peer-urls http://0.0.0.0:2380 \
   --initial-advertise-peer-urls http://0.0.0.0:2380 \
   --initial-cluster s1=http://0.0.0.0:2380 \
   --initial-cluster-token tkn \
   --initial-cluster-state new \
   --log-level info \
   --logger zap \
   --log-outputs stderr 
   `