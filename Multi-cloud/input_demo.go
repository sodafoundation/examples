package main

import (
	"context"
	"os"
	"fmt"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

const (
	bucket = "bkt001"
	KAFKA_URL = "localhost:9092"
	KAFKA_TOPIC = "multicloud-input"
	KAFKA_GROUP_ID = "TestGroupID"
)

var inps = [12] string {
	// Input DATA with labels

	// photo1.txt apple
	// photo2.txt zebra
	// photo3.txt parrot
	// photo4.txt orange
	// photo5.txt pear
	// photo6.txt lion
	// photo7.txt cat
	// photo8.txt dog
	// photo9.txt banana
	// photo10.txt duck
	// photo11.txt dove
	// photo12.txt mango

	"photo1.txt",
	"photo2.txt",
	"photo3.txt",
	"photo4.txt",
	"photo5.txt",
	"photo6.txt",
	"photo7.txt",
	"photo8.txt",
	"photo9.txt",
	"photo10.txt",
	"photo11.txt",
	"photo12.txt",

}

func getKafkaWriter(kafkaURL, topic, groupID string) *kafka.Writer {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func upload() {

	for _, fname := range inps  {
		fmt.Println("Uploading [", bucket, "] [", fname, "]")
		GelatoUpload(bucket, fname)
	}

}

func processFile(filename string) []byte {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error Opening input file", filename, err)
		return nil
	}
	
	defer f.Close()

	buf := make([]byte, 16)
    n, err := f.Read(buf)
	if err != nil {
		fmt.Println("Error Reading input file", filename, err)
		return nil
	}

	if n > 0 {
		// fmt.Println("[", n, "]", string(buf))
	}

	return buf
}

func main() {
	kafkaURL := KAFKA_URL
	topic := KAFKA_TOPIC
	groupID := KAFKA_GROUP_ID

	writer := getKafkaWriter(kafkaURL, topic, groupID)
	defer writer.Close()


	for _, fname := range inps  {
		fmt.Println("Downloading [" + bucket + "] [" + fname + "]")
		GelatoDownload(bucket, fname)

		value := string (processFile(fname))
		values := strings.Fields(value)
		value = values[0]

		fmt.Println("filename =", fname, "Value =", value)
		writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fname),
				Value: []byte(value),
			},
		)
	}
}
