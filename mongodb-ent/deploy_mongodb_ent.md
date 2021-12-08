
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

## Install a Replicaset Cluster with Prometheus Exporter
``` sh
kubectl apply -f ./secret2.yaml # create project in org-1 and add API key
kubectl apply -f ./replicaset-server.yaml # install sharded cluster with prometheus exporter
kubectl apply -f ./secret4.yaml # create database user for admin and prometheus exporter
```

## Install mongodb client and Connect to Mongodb
``` sh
kubectl apply -f ./mongo-client.yaml
kubectl exec -it $(kubectl get pod -n mongodb | grep "mongo-client" | awk '{print $1}') -n mongodb -- /bin/bash
mongo mongodb://admin:P%40ssw0rd@shard-svc:27017/admin?ssl=false # for shard
mongo mongodb://admin:P%40ssw0rd@replicaset-svc:27017/admin?ssl=false # for replicaset
```

## Insert Some Test Data and Do Sharding
``` sh
rs.slaveOk()  # run this when you are using replicaset mode and you are at secondary
use testdb
sh.enableSharding("testdb") # this can only be run on shard mode
sh.shardCollection("testdb.user", { _id : "hashed" } , false) # this can only be run on shard mode
db.user.insert({name: "kevin", age: "28"})
db.user.getShardDistribution() # this can only be run on shard mode
```

## Kubernetes View
* statefulset
``` sh
NAME             READY   AGE
ops-manager      1/1     7d2h
ops-manager-db   3/3     7d2h
replicaset       3/3     131m
shard-0          3/3     18m
shard-1          3/3     3m47s
shard-config     1/1     18m
shard-mongos     2/2     17m
```
* pod
``` sh
NAME                                           READY   STATUS    RESTARTS   AGE
mongo-client-847658c87d-b8c4m                  1/1     Running   0          7d
mongodb-enterprise-operator-6495bdd947-bvkg4   1/1     Running   0          7d2h
ops-manager-0                                  1/1     Running   0          7d2h
ops-manager-db-0                               3/3     Running   0          7d2h
ops-manager-db-1                               3/3     Running   0          7d2h
ops-manager-db-2                               3/3     Running   0          7d2h
replicaset-0                                   2/2     Running   0          130m
replicaset-1                                   2/2     Running   0          130m
replicaset-2                                   2/2     Running   0          130m
shard-0-0                                      2/2     Running   0          17m
shard-0-1                                      2/2     Running   0          17m
shard-0-2                                      2/2     Running   0          17m
shard-1-0                                      2/2     Running   0          3m17s
shard-1-1                                      2/2     Running   0          3m8s
shard-1-2                                      2/2     Running   0          2m56s
shard-config-0                                 1/1     Running   0          17m
shard-mongos-0                                 2/2     Running   0          16m
shard-mongos-1                                 2/2     Running   0          16m
``` 
* service
``` sh
NAME                      TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
operator-webhook          ClusterIP   10.102.41.150   <none>        443/TCP           7d2h
ops-manager-db-svc        ClusterIP   None            <none>        27017/TCP         7d2h
ops-manager-svc           ClusterIP   None            <none>        8080/TCP          7d2h
ops-manager-svc-ext       NodePort    10.103.74.52    <none>        8080:30036/TCP    7d2h
replicaset-svc            ClusterIP   None            <none>        27017/TCP         131m
replicaset-svc-external   NodePort    10.103.162.20   <none>        27017:30246/TCP   131m
shard-cs                  ClusterIP   None            <none>        27017/TCP         18m
shard-cs-external         NodePort    10.97.121.98    <none>        27017:31180/TCP   18m
shard-sh                  ClusterIP   None            <none>        27017/TCP         18m
shard-sh-external         NodePort    10.103.86.155   <none>        27017:32236/TCP   18m
shard-svc                 ClusterIP   None            <none>        27017/TCP         17m
shard-svc-external        NodePort    10.110.235.51   <none>        27017:31621/TCP   17m
```
* mongodb
``` sh
NAME         PHASE     VERSION     TYPE             AGE
replicaset   Running   4.4.0-ent   ReplicaSet       132m
shard        Running   4.4.0-ent   ShardedCluster   19m
```
* mongodbuser
``` sh
NAME                       PHASE     AGE
replicaset-admin           Updated   132m
replicaset-prom-exporter   Updated   132m
shard-admin                Updated   19m
shard-prom-exporter        Updated   19m
```
* configmap
``` sh
NAME                        DATA   AGE
kube-root-ca.crt            1      7d2h
ops-manager-db-project-id   1      5d22h
project-replicaset          3      141m
project-shard               3      6d1h
```


## Uninstall
``` sh
kubectl delete -f ./mongo-client.yaml
kubectl delete -f ./secret4.yaml
kubectl delete -f ./replicaset-server.yaml
kubectl delete -f ./secret3.yaml
kubectl delete -f ./shard-server.yaml
kubectl delete -f ./secret2.yaml
kubectl delete -f ./ops-manager.yaml
kubectl delete pvc --all -n mongodb
kubectl delete -f ./secret.yaml
kubectl delete -f ./mongodb-enterprise.yaml
kubectl delete -f ./crds.yaml
```

## Deploy simple mongodb producer to test mongodb HPA
* Use the helm in this repo: https://github.com/cwza/simple_mongodb
* It will deploy a producer to continuously generate read traffic to mongodb
``` sh
kubectl create namespace try
cd helm
b helm install --namespace=try -f values.yaml simple-mongod. # for shard
# helm install --namespace=try -f values.yaml --set mode=replicaset --set consumerurl="mongodb://admin:P%40ssw0rd@replicaset-svc.mongodb:27017/admin?ssl=false" simple-mongodb . # for replicaset
# helm delete simple-mongodb --namespace=try
```

## Prometheus Metrics
### Mongodb 0.1 Exporter
``` sh
# number of replicas
sum(mongodb_mongod_replset_member_health{pod=~"shard-0.*"})
kube_statefulset_status_replicas_ready{namespace="mongodb", statefulset="shard-0"}
# query rate (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"shard-[0-9].*"}[1m]))*60
# ThroughputThreshold (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"shard-[0-9].*"}[180m]))*60
# Avg. Read Latency (microseconds/queries)
sum(irate(mongodb_mongod_op_latencies_latency_total{type="read",namespace="mongodb",pod=~"shard-[0-9].*"}[1m]))/sum(irate(mongodb_mongod_op_latencies_ops_total{type="read",namespace="mongodb",pod=~"shard-[0-9].*"}[1m])) 
```
### Mongodb 0.2 Exporter
#### Shard
``` sh
# number of shards
mongodb_mongos_sharding_shards_total
# number of replicas
mongodb_mongod_replset_my_state # https://docs.mongodb.com/manual/reference/replica-states/
sum(mongodb_members_health{pod="shard-0-0", member_state=~"PRIMARY|SECONDARY"})
kube_statefulset_status_replicas_ready{namespace="mongodb", statefulset="shard-0"}
# query rate (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"shard-[0-9].*"}[1m]))*60
# ThroughputThreshold (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"shard-[0-9].*"}[180m]))*60
# Avg. Read Latency (microseconds/queries)
sum(irate(mongodb_ss_opLatencies_latency{op_type="reads",namespace="mongodb",pod=~"shard-[0-9].*"}[1m]))/sum(irate(mongodb_ss_opLatencies_ops{op_type="reads",namespace="mongodb",pod=~"shard-[0-9].*"}[1m])) 
```
#### Replicaset
``` sh
# number of replicas
mongodb_mongod_replset_my_state # https://docs.mongodb.com/manual/reference/replica-states/
sum(mongodb_members_health{pod="replicaset-0", member_state=~"PRIMARY|SECONDARY"})
kube_statefulset_status_replicas_ready{namespace="mongodb", statefulset="replicaset"}
# query rate (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"replicaset.*"}[1m]))*60
# ThroughputThreshold (queries/min)
sum(rate(mongodb_op_counters_total{type="query",namespace="mongodb",pod=~"replicaset.*"}[180m]))*60
# Avg. Read Latency (microseconds/queries)
sum(irate(mongodb_ss_opLatencies_latency{op_type="reads",namespace="mongodb",pod=~"replicaset.*"}[1m]))/sum(irate(mongodb_ss_opLatencies_ops{op_type="reads",namespace="mongodb",pod=~"replicaset.*"}[1m])) 
```