# RDS Operator
### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide in aws credentials here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart


# Features
- Create ( **All RDS Databases**, Subnet groups ) :white_check_mark:
- Delete ( **All RDS Databases**, Subnet groups ) :white_check_mark:
- Update ( COMING SOON ) :x:
- Create From RDS Snapshot ( **All RDS Databases** ) :white_check_mark:
- CreateReadReplica ( COMING SOON ) :x:

