package com.example.user_service.client;

import com.example.user_service.dto.EmailRequest;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class EmailKafkaProducer {

    private static final Logger logger = LoggerFactory.getLogger(EmailKafkaProducer.class);

    @Autowired
    private KafkaProducer<String, String> kafkaProducer;

    @Autowired
    private ObjectMapper objectMapper;

    @Value("${kafka.topic.email.requests:email_requests}")
    private String emailRequestTopic;

    public void queueEmailRequest(String to, String title, Map<String, String> data, String type) {
        try {
            String bodyJson = objectMapper.writeValueAsString(data);
            EmailRequest emailRequestPayload = new EmailRequest(to, title, bodyJson, type);
            
            // Serialize toàn bộ object EmailRequest thành JSON để gửi đi
            String finalPayload = objectMapper.writeValueAsString(emailRequestPayload);

            ProducerRecord<String, String> record = new ProducerRecord<>(emailRequestTopic, to, finalPayload);

            kafkaProducer.send(record, (metadata, exception) -> {
                if (exception == null) {
                    logger.info("Successfully queued email request for: {} to topic {}", to, metadata.topic());
                } else {
                    logger.error("Failed to send email request to Kafka for " + to, exception);
                }
            });

        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize email data to JSON for " + to, e);
        }
    }
}