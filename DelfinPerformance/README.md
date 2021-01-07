## Delfin Performance Monitoring with Prometheus & Kakfa

This is a demo of SODA Foundations Delfin performance monitoring feature using default prometheus exporter.

#### What is delfin
Delfin is a heterogeneous infrastructure management platform. It manages and monitors multiple storage backends from a single place. Prometheus is the default platform integrated with delfin for monitoring the performances of the storages.

#### What is prometheus

Prometheus is a free and popular open source system monitoring(event monitoring and alerting) tool. It uses a time-series database to store the real-time scpraped metrics value. The more detail about prometheus can be found [here](https://prometheus.io/).

#### How to connect prometheus with delfin

Below is the architecture diagram of prometheus and delfin. Delfin collects the performance metrics data through it's driver and pushes to the exporter interface. Prometheus exporter takes the data from the exporter interface of delfin and converted into time-series format(i.e prometheus database format) and persists in .txt file.

Delfin also run webserver client, which exposes this .txt file(metrics data) to http server on specific port.

Now, the prometheus server comes and scrapes the metrics over this targeted path.

![](/DelfinPerformance/delfin_architecture_with_prometheus.jpg)

#### What are the usecase
1. Users want to monitor and analyse the performance of storage-arrays
2. Users want to monitor and analyse the performance of storage-pools
3. Users want to monitor and analyse the performance of storage-volumes
4. Users want to monitor and analyse the performance of storage-controllers
5. Users want to monitor and analyse the performance of storage-ports
6. Users want to monitor and analyse the performance of storage-disks

The detailed usecases with different metrics(like read_throughput, bandwidth etc.) and resourcetype(like array, pool etc.) are [here](https://github.com/sodafoundation/design-specs/blob/dcdee7b67d4a4ee74f065f00b2e93efb22f2493a/specs/SIM/PerfomanceMontoringDesign.md)
#### Note:

  There are two exporters supported as of now for performance collection, prometheus and kafka. Uncomment performance_exporters option in <delfin_path>/etc/delffin/delfin.conf as below to enable exporter for prometheus and kafka  

  Ex:
   ```
# Uncomment or add exporters
performance_exporters = PerformanceExporterPrometheus, PerformanceExporterKafka
   ```

#### How to setup delfin with prometheus

  Follow the below steps to setup delfin with prometheus. Once your setup is ready, you can register the storage devices for performance monitoring. Later, the performance metrics can be viewed on prometheus server. This example also guides you to configure and update the targets and interval for scraping the metrics.

  Alternatively, you can also watch this [video]((https://drive.google.com/file/d/1WMmLXQeNlToZd0DP5hCFtDZ1IbNJpO6B/view?usp=drivesdk) for more detail.


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
            - localhost:8195
  ```

#### How to setup kafka with delfin

1. [Download](https://www.apache.org/dyn/closer.cgi?path=/kafka/2.6.0/kafka_2.13-2.6.0.tgz) the latest kafka release and extract it:
```
$ tar -xzf kafka_2.13-2.6.0.tgz
$ cd kafka_2.13-2.6.0
```

2. Start the kafka environment

  ```
  # Start the ZooKeeper service
  # Note: Soon, ZooKeeper will no longer be required by Apache Kafka.
  $ bin/zookeeper-server-start.sh config/zookeeper.properties
  ```
 Open another terminal session and run:

 ```
 # Start the Kafka broker service
$ bin/kafka-server-start.sh config/server.properties
 ```

3. Create a topic to publish and subscribe metrics data

  ```
  bin/kafka-topics.sh --create --topic delfin-kafka --bootstrap-server localhost:9092
  ```

4. run the kafka consumer client to see the collected metrics data

  kafka_consumer_client.py

  ```
  from kafka import KafkaConsumer
  import json

  def main():

    consumer = KafkaConsumer(
        'delfin-kafka',
        bootstrap_servers=['localhost:9092'],
        value_deserializer=lambda v: json.loads(v.decode('utf-8')))
    print(consumer.metrics())
    for message in consumer:
        print(message)

  if __name__ == '__main__':
    main()

  ```

### Installi Delfin and start performance collection


3. Follow this [link](https://github.com/sodafoundation/delfin/blob/master/installer/README.md) to install delfin

4. Register storage for performance collection

    Use API to register the storages for performance collection

    PUT http://localhost:8190/v1/storages/<storage_id>/metrics-config

    body:
    ```
    {
    "array_polling": {
    "perf_collection": true,
    "interval": 900,
    "is_historic": true
    }
    }
   ```
Example:

  ![](/DelfinPerformance/metri-config-api.png)



  Example1:

  ![](/DelfinPerformance/prometheus_dashboard.png)

  Example 2:

  ![](/DelfinPerformance/prometheus_dashboard2.png)

#### What user should see
  1. Performance metrics data on prometheus server
  2. The graphs of performances of storage devices

##### Note: Kafka exporter is supported from delfin release 1.1.0 onwards.  

