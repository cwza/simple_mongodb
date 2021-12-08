## Build and Run
``` sh
cd src
go build -o app
./app -cfgpath=./config.toml
```

## Deploy to Dockerhub
When you push to master branch the github action will automatically build image and push it to my dockerhub

## Deploy to k8s
``` sh
kubectl create namespace try
cd helm
b helm install --namespace=try -f values.yaml simple-mongod. # for shard
# helm install --namespace=try -f values.yaml --set mode=replicaset --set consumerurl="mongodb://admin:P%40ssw0rd@replicaset-svc.mongodb:27017/admin?ssl=false" simple-mongodb . # for replicaset
# helm delete simple-mongodb --namespace=try
```

## Default Value of Mongodb Driver
```
ConnectTimeout         30 * time.Second
Compressors            nil (compression will not be used)
Dialer                 net.Dialer with a 300 second keepalive time
HeartbeatInterval      10 * time.Second
LocalThreshold         15 * time.Millisecond
MaxConnIdleTime        nil (no limit)
MaxPoolSize            100
Monitor                nil
ReadConcern            nil (server default `local`)
ReadPreference         readpref.Primary()
Registry               bson.DefaultRegistry
RetryWrites            true
ServerSelectionTimeout 30 * time.Second
Direct                 false
SocketTimeout          nil (infinite)
TLSConfig              nil
WriteConcern           nil (server default `w:1`)
ZlibLevel              6 (if zlib compression enabled)
```

---

## Install Mongodb Enterprised in K8S
* Please see ./mongodb-ent/deploy_mongodb_ent.md