# RDS Operator

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/733bc755125b4a5193016e1fbf460700)](https://app.codacy.com/app/agill17/rds-operator?utm_source=github.com&utm_medium=referral&utm_content=agill17/rds-operator&utm_campaign=Badge_Grade_Settings)

### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide aws credentials in here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart


#### Note: If you want to deploy mysql and or postgress using aurora, you MUST use both DBCluster and DBInstance

# Features
- _**Create**_ ( **All RDS Databases**, Subnet groups )
  - DBCluster: `spec.createClusterSpec` :white_check_mark:
  - DBInstance: `spec.createInstanceSpec` :white_check_mark: 
- _**Delete**_ ( **All RDS Databases**, Subnet groups ) :white_check_mark:
  - DBCluster: `spec.deleteClusterSpec` :white_check_mark:
  - DBInstance: `spec.deleteInstanceSpec` :white_check_mark: 
- _**Restore**_ from rds snapshot ( **All RDS Databases** )
  - DBCluster: `spec.createClusterFromSnapshot` :white_check_mark:
  - DBInstance: `spec.createInstanceFromSnapshot` :white_check_mark:

## TODO 
- Reconcile AWS resources periodically :x: ( Recreate DB incase they get deleted )
- Update ( COMING SOON ) :x:
- CreateReadReplica when deployed from scratch and or from snapshot ( COMING SOON ) :x:

