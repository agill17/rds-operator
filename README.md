# RDS Operator

### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide aws credentials in here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart

---
#### Example/Samples
- Can be found under `examples/` dir
#### DBCluster vs DBInstance
- Use DBCluster and DBInstance (Both) **WHEN** creating anything with aurora ( so aurora-mysql or aurora-postgresql )
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
- Once DBCluster is created, you must attach instances to that cluster, thats where DBInstance comes into the picture.
- How to attach a DBInstance to a DBCluster? You have to provide DBClusterIdentifier within DBInstance.
- **WHEN** to use standalone DBInstance? -- **WHEN**ever you want to create anything thats __*not*__ aurora.

# Features
- _**Create**_ ( **All RDS Databases**, Subnet groups )
  - DBCluster: `spec.createClusterSpec` :white_check_mark:
    - credentialsFromSecret: `spec.credentialsFrom` :white_check_mark:
    - externalNameService: `spec.serviceName` :white_check_mark:
    - initDBJob `spec.initDBJob` :x:
  - DBInstance: `spec.createInstanceSpec` :white_check_mark: 
    - credentialsFromSecret: `spec.credentialsFrom` :x:
    - externalNameService: `spec.serviceName` :white_check_mark: ( primary db instance )
    - initDBJob `spec.initDBJob` :x: ( if attached to DBCluster, this wont run at all. )
- _**Delete**_ ( **All RDS Databases**, Subnet groups ) :white_check_mark:
  - DBCluster: `spec.deleteClusterSpec` :white_check_mark:
  - DBInstance: `spec.deleteInstanceSpec` :white_check_mark: 
- _**Restore**_ from rds snapshot ( **All RDS Databases** )
  - DBCluster: `spec.createClusterFromSnapshot` :white_check_mark:
  - DBInstance: `spec.createInstanceFromSnapshot` :white_check_mark:


## TODO 
~~- Centralize (create/delete/restore) into rdsLib and make use of existing interface funcs!~~ :white_check_mark:
- Add docs
- Make Secret resource conditional for both DBCluster(:white_check_mark:) and DBInstance ( some folks might want to pass credentials like username and password from a k8s secret so there is no need to deploy another secret with the same information )
- Add support for initDB Job in dbCluster so that a user can have their db imported from an existing image ( DBCluster & DBInstance )
- Make delete optional and get rid of deleteSpec to enforce snapshotting before deleting rds.
- Reconcile AWS resources periodically :x: ( Recreate DB incase they get deleted )
- Update ( COMING SOON ) :x:
- CreateReadReplica **WHEN** deployed from scratch and or from snapshot ( COMING SOON ) :x:

