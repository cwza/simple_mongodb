
## Install Operator and Ops-Manager
* https://github.com/mongodb/mongodb-enterprise-kubernetes
* https://github.com/mongodb/mongodb-enterprise-kubernetes/tree/master/samples
* https://docs.mongodb.com/manual/reference/operator/
``` sh
kubectl create namespace mongodb
kubectl apply -f ./crds.yaml # add crds
kubectl apply -f ./mongodb-enterprise.yaml # install operator
kubectl apply -f ./secret1.yaml # account for ops-manager and database user
kubectl apply -f ./ops-manager.yaml # install ops-manager
```

## Use Ops-Manager to Add Organization and API Key
* find the nodeip and external port from ops-manager-svc-ext service
* Go to Ops-Manager GUI, http://nodeip:30036
* Use admin/P@ssw0rd to login
* organizations > add organization named org-1
* Settings > Organization ID > Copy that id to ./secret2.yaml > orgId
* select the org-1 > access manager > API Keys > Create an API Key with Organization Owner permission > add the operator ip into api access list
* Copy the public key and private key into ./secret2.yaml > publicKey, privateKey

## Install a Sharded Cluster with Prometheus Exporter
``` sh
kubectl apply -f ./secret2.yaml # create project in org-1 and add API key
kubectl apply -f ./shard-server.yaml # install sharded cluster with prometheus exporter
kubectl apply -f ./secret3.yaml # create database user for admin and prometheus exporter
```

## Install mongodb client and Connect to Mongodb
``` sh
kubectl apply -f ./mongo-client.yaml
kubectl exec -it $(kubectl get pod -n mongodb | grep "mongo-client" | awk '{print $1}') -n mongodb -- /bin/bash
mongo mongodb://admin:P%40ssw0rd@shard-svc:27017/admin?ssl=false
```

## Insert Some Test Data and Do Sharding
```
use testdb
sh.enableSharding("testdb")
sh.shardCollection("testdb.user", { _id_ : "hashed" } , false)
db.user.insert({name: "kevin", age: "28"})
db.user.getShardDistribution()
```

## Uninstall
``` sh
kubectl delete -f ./mongo-client.yaml
kubectl delete -f ./secret3.yaml
kubectl delete -f ./sharded1.yaml
kubectl delete -f ./secret2.yaml
kubectl delete -f ./ops-manager.yaml
kubectl delete pvc --all -n mongodb
kubectl delete -f ./secret.yaml
kubectl delete -f ./mongodb-enterprise.yaml
kubectl delete -f ./crds.yaml
```
