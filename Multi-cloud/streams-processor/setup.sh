#!/bin/bash

# set -x


INPUT_TOPIC="multicloud-input"
OUTPUT_TOPIC="multicloud-output"
INPUT2_TOPIC="multicloud-input2"
OUTPUT2_TOPIC="multicloud-output2"

KAFKA_VER="kafka_2.12-2.4.1"

TOP_DIR=$(cd $(dirname "$0") && cd .. && pwd)
KAFKA_FOLDER=$TOP_DIR/$KAFKA_VER

echo "$TOP_DIR"
echo "$KAFKA_FOLDER"

# usage function
function usage()
{
   cat << HEREDOC

   Usage: setup.sh [ options ]

   optional arguments:
     -h, --help           show this help message and exit
     -c                   Clean and purge current installation
     -d                   Download kafka
     -s                   Start kafka

HEREDOC
}  


function download_kafka()
{
    echo "Download Kafka"
    cd $TOP_DIR
    wget https://downloads.apache.org/kafka/2.4.1/kafka_2.12-2.4.1.tgz
    tar -xzf kafka_2.12-2.4.1.tgz
}

function start_kafka()
{
    cd $KAFKA_FOLDER

    # Update Kafka config for delete topic support, if needed
    grep -qxF 'delete.topic.enable=true' config/server.properties || echo 'delete.topic.enable=true' >> config/server.properties

    echo "Start Zookeeper"
    bin/zookeeper-server-start.sh config/zookeeper.properties >zookeeper.log 2>&1 &
    sleep 2
    


    echo "Start Kafka"
    bin/kafka-server-start.sh config/server.properties >kafka.log 2>&1 &
    sleep 2

    # echo "Create Kafka topics"
    # bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic $INPUT_TOPIC
    # bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic $OUTPUT_TOPIC
    # bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic $INPUT2_TOPIC
    # bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic $OUTPUT2_TOPIC

    # echo "Kafka console"
    # bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic $INPUT_TOPIC
    # bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic $OUTPUT2_TOPIC

    echo "List Topics"
    bin/kafka-topics.sh --list --bootstrap-server localhost:9092
}

function clean()
{
    cd $KAFKA_FOLDER

    # echo "Remove topics of Kafka"
    # bin/kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic $INPUT_TOPIC
    # bin/kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic $OUTPUT_TOPIC
    # bin/kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic $INPUT2_TOPIC
    # bin/kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic $OUTPUT2_TOPIC

    echo "Stopping Kafka"
    bin/kafka-server-stop.sh

    echo "Stopping Zookeeper"
    bin/zookeeper-server-stop.sh
}

function start_app()
{
    cd $KAFKA_FOLDER

    echo "Start Stream Processor"
    mvn clean package
    mvn exec:java -Dexec.mainClass=myapps.DemoMultiCloudInput
    # mvn exec:java -Dexec.mainClass=myapps.DemoMultiCloudOutput
}

OPTIND=1         # Reset in case getopts has been used previously in the shell.

while getopts ":h?cds" opt; do
    case "$opt" in
    h|\?)
        usage
        exit 0
        ;;
    :)  echo "option -$OPTARG requires an argumnet"
        usage
        ;;
    c)  clean="1"
        ;;
    d)  download="1"
        ;;
    s)  startkafka="1"
        ;;

    esac
done

shift $((OPTIND-1))

[ "${1:-}" = "--" ] && shift

if [ "$clean" = "1" ]; then
    clean
    exit
fi

if [ "$download" = "1" ]; then
    download_kafka
    exit
fi

if [ "$startkafka" = "1" ]; then
    start_kafka
    exit
fi

# --- Start ---

download_kafka
start_kafka
# start_app
