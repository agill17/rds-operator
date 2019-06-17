# RDS Operator

### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide aws credentials in here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart

---
#### DBCluster vs DBInstance
- Use DBCluster when creating anything with aurora ( so aurora-mysql or aurora-postgresql )
- Why use DBCluster? Because of how AWS SDK works.
- Their docs for DBCluster state the following
  ```The name of the database engine to be used for this DB cluster.
    //
    // Valid Values: aurora (for MySQL 5.6-compatible Aurora), aurora-mysql (for
    // MySQL 5.7-compatible Aurora), and aurora-postgresql
    //
    // Engine is a required field
    Engine *string `type:"string" required:"true"`
   ```
- So the only valid values for DBClusters are aurora-mysql and aurora-postgresql 

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
- Add docs
- Make Secret resource conditional for both DBCluster and DBInstance ( some folks might want to pass credentials like username and password from a k8s secret so there is no need to deploy another secret with the same information )
- Add support for initDB Job in dbCluster so that a user can have their db imported from an existing image ( DBCluster & DBInstance )
- Make delete optional and get rid of deleteSpec to enforce snapshotting before deleting rds.
- Reconcile AWS resources periodically :x: ( Recreate DB incase they get deleted )
- Update ( COMING SOON ) :x:
- CreateReadReplica when deployed from scratch and or from snapshot ( COMING SOON ) :x:

