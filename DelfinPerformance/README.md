## Delfin Performance Monitoring with Prometheus

This is demo of SODA Foundations Delfin performance monitoring feature using default preometheus exporter.


Delfin is a heterogeneous infrastructure management platform. It manages and monitors multiple storage backends from a single place. Prometheus is default platform integrated with delfin for monitoring the performances of the storages.

This example helps to setup delfin and prometheus together for POC. Once your setup is ready, you can register the storage devices for performance monitoring. Later, the performance metrics can be viewed on prometheus server. This example also guides you to configure and update the targets and interval for scraping the metrics. 


Follow the below steps to setup delfin and prometheus or go through the demo [video](https://drive.google.com/file/d/1WMmLXQeNlToZd0DP5hCFtDZ1IbNJpO6B/view?usp=drivesdk) for the reference.

step1: Install and start prometheus server

1. [Download the latest binaries from here](https://prometheus.io/download/) and run the below steps.

  ```sh
  1. tar xvfz prometheus-*.tar.gz

  2. cd prometheus-*

  3. ./prometheus
  ```
Example:
```sh
root@root:/prometheus/prometheus-2.20.0.linux-amd64# ./prometheus
```

2. Edit the prometheus.yml and set the appropriate target, interval and metrics_api   path. Below is sample example of prometheus.yml

  ###### prometheus.yml

  ```
  global:
    scrape_interval: 10s
  scrape_configs:
   - job_name: example
   metrics_path: /metrics
   static_configs:
    - targets:
            - localhost:5000
  ```

3. Install and run the Delfin from [here](https://github.com/sodafoundation/delfin/blob/master/installer/README.md)


4. Register and update performance metric collection(as of now only array metrics collection is supported) by below API.
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

5. Run below client.py program to start webserver(it exposes the metrics to https server)

  ##### client.py

  ```
  from flask import Flask

  app = Flask(__name__)

  @app.route("/metrics", methods=['GET'])
  def getfile():
    with open("/var/lib/delfin/delfin_exporter.txt", "r+") as f:
        data = f.read()
    return data

  if __name__ == '__main__':
      app.run(host='localhost')
  ```

6. Now, the collected metrics can be seen on prometheus server

  Example1:

  ![](/DelfinPerformance/prometheus_dashboard.png)

  Example 2:

  ![](/DelfinPerformance/prometheus_dashboard2.png)
