# RDS Operator
### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide aws credentials in here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart


# Features
- Create ( **All RDS Databases**, Subnet groups ) :white_check_mark:
- Delete ( **All RDS Databases**, Subnet groups ) :white_check_mark:
- Create from an RDS Snapshot ( **All RDS Databases** ) :white_check_mark:
- Reconcile AWS resources periodically :x: ( Recreate DB incase they get deleted )
- Update ( COMING SOON ) :x:
- CreateReadReplica ( COMING SOON ) :x:

