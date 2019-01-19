
# RDS-Operator
# Installation
##### Note: The operator requires AWS Creds to be passed in as env variables ( could be passed from secrets )

1.  *Operator.yaml*
    ```
    $ kubectl apply -f example/operator.yaml
    ```
2. *RDS CR (MySQL)*
    ```
    $ kubectl apply -f example/mysql_instance_sample_cr.yaml
    ```
    *OR RDS CR (MySQL)  (Aurora MySQL)*
    ```
    $ kubectl apply -f example/aurora_mysql_sample_cr.yaml
    ```
# Features;
  - RDS DB Provisioning
  - Pre-Populate provisioned DB ( *optional* )
  - Reheal from latest available snapshot when DB no longer exists in AWS ( *optional* )
# Supported DB Engines
- MySQL
- Aurora-MySQL
# Todos
 - Cleanup
 - Add DBSubnetGroup kind
 - Support more db engines


