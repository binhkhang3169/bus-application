package com.example.trip_service.config;

import java.util.Properties;

import org.apache.kafka.clients.CommonClientConfigs;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.config.SaslConfigs;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.StringSerializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class KafkaClientConfig {

    // Đọc các giá trị từ application.properties
    @Value("${spring.kafka.bootstrap-servers}")
    private String bootstrapServers;

    @Value("${spring.kafka.security.protocol}")
    private String securityProtocol;

    @Value("${spring.kafka.sasl.mechanism}")
    private String saslMechanism;

    @Value("${spring.kafka.sasl.username}")
    private String saslUsername;

    @Value("${spring.kafka.sasl.password}")
    private String saslPassword;

    /**
     * Phương thức private để tạo cấu hình cơ sở, bao gồm cả bảo mật.
     * @return Properties chứa cấu hình SASL/SSL.
     */
    private Properties createBaseProperties() {
        Properties props = new Properties();
        props.put(CommonClientConfigs.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        props.put(CommonClientConfigs.SECURITY_PROTOCOL_CONFIG, securityProtocol);
        props.put(SaslConfigs.SASL_MECHANISM, saslMechanism);
        
        // Tạo chuỗi JAAS config để xác thực username và password
        String jaasConfig = String.format(
                "org.apache.kafka.common.security.scram.ScramLoginModule required username=\"%s\" password=\"%s\";",
                saslUsername, saslPassword
        );
        props.put(SaslConfigs.SASL_JAAS_CONFIG, jaasConfig);
        
        return props;
    }

    /**
     * Tạo một bean KafkaProducer đã được cấu hình bảo mật đầy đủ.
     */
    @Bean
    public KafkaProducer<String, String> kafkaProducer() {
        Properties props = createBaseProperties(); // Lấy cấu hình bảo mật
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        return new KafkaProducer<>(props);
    }

    /**
     * Phương thức tạo một instance KafkaConsumer mới với cấu hình bảo mật.
     */
    public KafkaConsumer<String, String> createKafkaConsumer(String groupId) {
        Properties props = createBaseProperties(); // Lấy cấu hình bảo mật
        props.put(ConsumerConfig.GROUP_ID_CONFIG, groupId);
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");
        return new KafkaConsumer<>(props);
    }

    // @Bean
    // public ObjectMapper objectMapper() {
    //     return new ObjectMapper();
    // }
}