package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	KafkaBootstrapServers string
	KafkaUsername         string
	KafkaPassword         string
	KafkaSaslMechanism    string
	KafkaSecurityProtocol string
	KafkaGroupId          string
	KafkaTopic            string
	KafkaTopicDlq         string
	TargetServiceUrl      string
)

func Init() error {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading configuration file: %v", err)
	}

	KafkaBootstrapServers = viper.GetString("KAFKA_BOOTSTRAP_SERVERS")
	KafkaUsername = viper.GetString("KAFKA_USERNAME")
	KafkaPassword = viper.GetString("KAFKA_PASSWORD")
	KafkaSaslMechanism = viper.GetString("KAFKA_SASL_MECHANISM")
	KafkaSecurityProtocol = viper.GetString("KAFKA_SECURITY_PROTOCOL")
	KafkaGroupId = viper.GetString("KAFKA_GROUP_ID")
	KafkaTopic = viper.GetString("KAFKA_TOPIC")
	KafkaTopicDlq = viper.GetString("KAFKA_TOPIC_DLQ")
	TargetServiceUrl = viper.GetString("TARGET_SERVICE_URL")

	log.Info("Configuration loaded successfully")

	return nil
}
