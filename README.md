# RDS Operator

### Install and Deploy 

#### Option 1: kubectl way
1. `kubectl apply -f examples/operator.yaml` ( provide aws credentials in here )
2. `kubectl apply -f examples/auroraMySQLCluster.yaml`

#### Option 2: Helm chart

---
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
  - DBInstance: `spec.createInstanceSpec` :white_check_mark: 
- _**Delete**_ ( **All RDS Databases**, Subnet groups ) :white_check_mark:
  - DBCluster: `spec.deleteClusterSpec` :white_check_mark:
  - DBInstance: `spec.deleteInstanceSpec` :white_check_mark: 
- _**Restore**_ from rds snapshot ( **All RDS Databases** )
  - DBCluster: `spec.createClusterFromSnapshot` :white_check_mark:
  - DBInstance: `spec.createInstanceFromSnapshot` :white_check_mark:

### What really happens **WHEN** cr's are created for DBCluster 
---
1. Validates required input
2. Attempts to create a new secret 
    
    - Secret __does not get created__
        -  **WHEN** spec.credentialsFrom __is defined__
        grab usernameKey, passwordKey and secretName and store into cr.status ( just the keys )
    - Secret __gets created__ 
      - **WHEN** spec.credentialsFrom __is NOT defined__ and spec is of type createClusterSpec ( not createClusterFromSnapshot )
        - If spec.createClusterSpec.MasterUsername && cr.spec.createClusterSpec.MasterUserPassword __is defined__, secret gets created, finally the secretName, userKey and passKey gets stored in CR status ( just the keys )
          ```
          secretName: cr.Name-secret
          data:
            DB_USER: spec.createClusterSpec.MasterUsername
            DB_PASS: cr.spec.createClusterSpec.MasterUserPassword
          ```
        - If spec.createClusterSpec.MasterUsername && cr.spec.createClusterSpec.MasterUserPassword __is NOT defined__, secret gets created, finally the secretName, userKey and passKey gets stored in CR status ( just the keys ) 
          ```
          secretName: cr.Name-secret
          data:
            DB_USER: admin
            DB_PASS: password
          ```

3. Updates fields in CRs
4. Checks if rds is already created ( using status fields of CR )
5. If not, checks for existence and creates the rds cluster object in AWS
6. Waits until cluster is available
7. Once available, update CR status ( created: true )
9. Create a external svc and point to cluster endpoint

- DBInstance CR; 
  1. Validates required input
  2. sets up values of fields that can be assumed
  3. If part of cluster ( meaning clusterID is specified )
    - Checks and waits until cluster is available by quering aws for status
  4. Checks if rds is already created ( using status field of the CR )
  5. If no, checks for existence and creates the rds cluster object in AWS
  6. Waits until instance is available
  7. Once available, update statuses ( created: true )
  8. Creates and deployes secret ( optional )
  9. Creates external name svc with instance endpoint ( only if not part of cluster  )
  10. Repeat


## TODO 
- Add docs
- Make Secret resource conditional for both DBCluster and DBInstance ( some folks might want to pass credentials like username and password from a k8s secret so there is no need to deploy another secret with the same information )
- Add support for initDB Job in dbCluster so that a user can have their db imported from an existing image ( DBCluster & DBInstance )
- Make delete optional and get rid of deleteSpec to enforce snapshotting before deleting rds.
- Reconcile AWS resources periodically :x: ( Recreate DB incase they get deleted )
- Update ( COMING SOON ) :x:
- CreateReadReplica **WHEN** deployed from scratch and or from snapshot ( COMING SOON ) :x:

