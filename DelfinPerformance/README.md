## Delfin Performance Monitoring with Preometheus

This is demo of SODA Foundations Delfin performance monitoring feature using default preometheus exporter.


Delfin is a heterogeneous infrastructure management platform. It manages and monitors multiple storage backends from a single place. Prometheus is default integrated with delfin for monitoring performances of the storages.

The detailed analysis and architecture of prometheus integration with delfin can be found [here](https://github.com/sodafoundation/design-specs/blob/dcdee7b67d4a4ee74f065f00b2e93efb22f2493a/specs/SIM/PerfomanceMontoringDesign.md)

Follow the below steps to setup delfin and prometheus or go through the demo [video]() for the reference.

step1: Install and start prometheus server

1. [Download the latest binaries from here](https://prometheus.io/download/) and run the below steps.

  ```sh
  1. tar xvfz prometheus-*.tar.gz

  2. cd prometheus-*

  3. ./prometheus
  ```

2. Example:
```sh
root@root:/prometheus/prometheus-2.20.0.linux-amd64# ./prometheus
```

3. Install and run the Delfin from [here](https://docs.sodafoundation.io/soda-gettingstarted/installation-using-ansible/)

  make sure, delfin is enabled in configuration file. i.e Update the file ansible/group_vars/delfin.yml and change the value of enable_delfin to true

  Example:

  ##### vim ansible/group_vars/delfin.yml
  ```
  # Install delfin (true/false)
  enable_delfin: true
  ```

4. Update the metrics variable(ex, interval, is_historic). Currently "array" metrics type is supported for performance monitoring. This can be updated bu using belwo API:
```
http://localhost:8190/v1/storages/storage_id/metrics-config
```
The body can be used to update the variables:
```
{
  "array_polling": {
    "perf_collection": true,
    "interval": 6,
    "is_historic": true
    }
}
```
Example:

  ![](/DelfinPerformance/metri-config-api.png)

5. Now, the collected metrics can be seen on prometheus server

Example1:

![](/DelfinPerformance/prometheus_dashboard.png)

Example 2:

![](/DelfinPerformance/prometheus_dashboard2.png)
