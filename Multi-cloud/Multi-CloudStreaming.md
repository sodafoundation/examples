# Multi-Cloud Streaming

This Project contains demos that shows how to use SODA Foundations multi-cloud project for streaming usecases. The Apache Kafka project is used as the Stream Processing platform for these use demos.

Two cases are identified for the purpose of demonstrating multi-cloud features for Streaming. We can connect the Stream Processor in two ways to multi-cloud as described below.

## Case 1: Connecting Stream Processor after Multi-Cloud

In this case Stream Processor get its input from Multi-cloud, where this inputs are objects stored in cloud and/or on-prem. These objects can be previously stored Stream Processor output, or any objects that are stored and managed by multi-cloud.

This demo shows how to get data/files from the cloud using multi-cloud, and do stream processing like analysis & categorization and generate processed output. In this simplified usecase, the Stream processor download a list of photo files, which are labelled with names of animals, birds and fruits. Stream processor will identify these from labels and categorize them to Animals, Birds, Fruits and Unknown and outputs, the count of these categories.

 ### Steps to setup this demo:
  - Deploy multi-cloud v0.10.0 with Dashboard and upload input data
  - Configure and Start Kafka streaming
  - Start Kafka console client for output checking 
  - Start Kafka demo application that process on input stream data
  - Start SODA multi-cloud services
  - Start GO application that uses multi-cloud client to download data for iinput to stream processor

 ### Detailed steps
 - Deploy SODA Multi-Cloud and Dashboard Project & upload input data
Refer document https://github.com/sodafoundation/opensds/wiki/OpenSDS-Cluster-Installation-through-Ansible0

    ```bash
    git clone https://github.com/sodafoundation/opensds-installer.git
    cd opensds-installer
    git checkout v0.10.0

    # change  "opensds-installer/ansible/group_vars/common.yml
    host_ip: 192.168.20.128
    deploy_project: gelato

    # Deploy multi-cloud and keystone and Dashboard
    cd ansible
    chmod +x ./install_ansible.sh && ./install_ansible.sh
    ansible-playbook site.yml -i local.hosts
    ```

    Use Dashboard GUI to create an AKSK, add backend, create bucket (`bkt001`) and upload input files for demo

 - Start Kafka, Zookeeper,  and create topics
    ```bash
    git clone https://github.com/sodafoundation/demos
    cd ./demos/mc-demo/streams-processor

    ./setup.sh  # Download and Start Kafka, Create kafka topics
   cd kafka_2.12-2.4.1/
   bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic multicloud-input
   bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic multicloud-input2
   bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic multicloud-output
   bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic multicloud-output2
    ```

 - Using SODA Dashboard upload input data files in folder ./mc-demo/DATA
    ```bash
    # DATA folder contains 12 files with photo*.txt to be uploaded 
    ```

 - Build and start JAVA application that do the Stream processing
    ```bash
    cd ./mc-demo/streams-processor
    mvn clean package
    mvn exec:java -Dexec.mainClass=myapps.DemoMultiCloudInput
    ```

 - Start Kafka console client for checking the generated output
    ```bash
    # Open a new terminal
    cd ./mc-demo/kafka_2.12-2.4.1
    bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic multicloud-output

    ```

 - Start GO app with multi-cloud client
    ```bash
    # Untar multi-cloud branch with Client support
    # Open Another terminal, and

    mkdir -p  <gopath>/src/github.com/opensds
    tar xzvf multi-cloud.tar.gz -C <gopath>/src/github.com/opensds 
    cd <gopath>/src/github.com/opensds/multi-cloud
    make all
    
    #Get AKSK generated from OpenSDS Dashboard UI
    export OS_ACCESS_KEY=5gWcxQ9HfzJ1xkvmQMz6

    # Export environment variables
    export MULTI_CLOUD_IP=192.168.20.158
    export MICRO_SERVER_ADDRESS=:8089
    export OS_AUTH_AUTHSTRATEGY=keystone

    export OPENSDS_ENDPOINT=http://192.168.20.158:50040
    export OPENSDS_AUTH_STRATEGY=keystone
    export OS_AUTH_URL=http://192.168.20.158/identity
    export OS_USERNAME=admin
    export OS_PASSWORD=opensds@123
    export OS_TENANT_NAME=admin
    export OS_PROJECT_NAME=admin
    export OS_USER_DOMIN_ID=default

    export KAFKA_URL=localhost:9092
    export KAFKA_TOPIC=streams-wordcount-processor-output
    export KAFKA_GROUP_ID=TestGroupID
    
    # Run mc-demo project
    cd /path/to/mc-demo
    go run input-demo.go gelato_client.go

    ```


 - Check the Java APP output and output topic
    ```bash
   # Check JAVA app output & output topic
    ```
 
 ### Case 1 Demo Video
 [Demo1.mp4](https://drive.google.com/open?id=1J1pNPLuyxi9oIj9YzzD-WRkJfdBrwbMh)


## Case 2: Connecting Stream processing in front of Multi-Cloud
In this case we convert stream output from the Stream processor and create an object that multicloud can work on. We use Apache Kafka as Stream Processing platform.

Scenario for this kind of use-case is, when we need to store a report generated by a streaming framework into multiple cloud back ends for later consumption. We use multi-cloud services of Data Mover or Migration to store these objects in the Multi-cloud backend.

This demo shows a stream processing application which counts words on the console input, and generate a word count table. A go application will use this table to generate a file and store the file object to a cloud storage and apply different policies on the uploaded file using multi-cloud project.

### Steps to setup this demo:
  - Deploy multi-cloud v0.10.0 with Dashboard
  - Configure and Start Kafka streaming
  - Start Kafka console producer for input generation
  - Start Kafka demo application that process on input stream data
  - Start SODA multi-cloud services
  - Start GO application that uses multi-cloud client to upload output of stream processor

 ### Detailed steps
  - Deploy multi-cloud with Dashboard
    ```bash
    # Similar to Demo1, but no need to updload input files.
    ```
  - Configure and Start kafka, and create topics streaming
    ```bash
    # Similar to Demo1
    ```
  - Start Kafka console producer for input capturing for processing
    ```bash
    #Start in a new terminal window
    bin/kafka-console-producer.sh --broker-list localhost:9092 --topic multicloud-input2

    >this is test
    >another test
    ```

  - Start Kafka demo application that process on input stream data
    ```bash
    cd ./mc-demo/streams-processor
    mvn clean package
    mvn exec:java -Dexec.mainClass=myapps.DemoMultiCloudOutput
    ```

  - Start GO application that uses multi-cloud client
    ```bash
    cd ./mc-demo
    go run input-demo.go gelato_client.go

    ```
  - Verify generated output file from Dashboard for the count output
    ```bash
    # Dashboard will contain a file out.txt, with output
    ```

 ### Case 2 Dem2 Video
 [Demo2.mp4](https://drive.google.com/open?id=1l01C20C4eQypNw-rKeN8CEk68GusII1P)


