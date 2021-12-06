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
helm install --namespace=try -f values.yaml simple-mongodb .
helm delete simple-mongodb --namespace=try
```

---

## Install Mongodb Enterprised in K8S
* Please see ./mongodb-ent/deploy_mongodb_ent.md