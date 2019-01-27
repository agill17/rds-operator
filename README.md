
# RDS-Operator
# Installation
##### Note: The operator requires AWS Creds to be passed in as env variables ( could be passed from secrets ) 
##### OR 
##### the operator pod must be ran on a node that has an IAM role with priviliges to RDS attached.

1.  *Operator.yaml* ( AWS Creds go in here as env vars )
    ```
    $ kubectl apply -f example/operator.yaml
    ```
2. *Install CR for MySQL*
    ```
    $ kubectl apply -f example/mysql_instance_sample_cr.yaml
    ```
    *OR CR for Aurora MySQL*
    ```
    $ kubectl apply -f example/aurora_mysql_sample_cr.yaml
    ```
3. *Get Installed CR's*
    ```
    $ kubectl get dbInstances --all-namespaces
    $ kubectl get dbClusters --all-namespaces
# Features;
  - RDS DB Provisioning ( Instance + Cluster -- Separate kinds for each )
  - RDS DB Subnet Group Creation ( Separate kind )
  - Pre-Populate/Initialize provisioned DB ( *optional* )
  - Reheal from latest available snapshot when DB no longer exists in AWS ( *optional* )
  - Cleans up AWS resources when CR/Namespace is deleted
# Supported DB Engines
- MySQL
- Aurora-MySQL
# TODOS
 - Add DBParameterGroup controller
 - Paramterized periodic existence check per controller ( AWS Resources )
 - Handle modifications of a deployed CR's ( Update )
 - Support more db engines


