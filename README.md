
# RDS-Operator
# Installation
##### Note: The operator requires AWS Creds to be passed in as env variables ( could be passed from secrets )

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
  - Pre-Populate provisioned DB ( *optional* )
  - Reheal from latest available snapshot when DB no longer exists in AWS ( *optional* )
  - Cleans up AWS resources when CR/Namespace is deleted
# Supported DB Engines
- MySQL
- Aurora-MySQL
# Todos
 - Add DBSParameterGroup controller
 - Handle modifications of CR's
 - Support more db engines


