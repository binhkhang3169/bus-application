package com.example.trip_service.service;

import com.example.trip_service.config.KafkaClientConfig;
import com.example.trip_service.dto.TripStatusUpdateEvent;
import com.example.trip_service.model.Trip;
import com.example.trip_service.service.intef.TripService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.messaging.simp.SimpMessageSendingOperations;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicBoolean;

// Thay @Service bằng @Component và bỏ @KafkaListener
@Component
public class TripStatusListener {

    private static final Logger logger = LoggerFactory.getLogger(TripStatusListener.class);

    private final SimpMessageSendingOperations messagingTemplate;
    private final ObjectMapper objectMapper;
    private final TripService tripService;
    
    // === BẮT ĐẦU: Thêm các thành phần quản lý thủ công ===
    private final KafkaConsumer<String, String> consumer;
    private final ExecutorService executorService = Executors.newSingleThreadExecutor();
    private final AtomicBoolean running = new AtomicBoolean(false);

    private static final String TOPIC = "trip_status_updated";
    private static final String GROUP_ID = "trip_status_websocket_group";
    // === KẾT THÚC ===

    @Autowired
    public TripStatusListener(SimpMessageSendingOperations messagingTemplate, ObjectMapper objectMapper, TripService tripService, 
                              KafkaClientConfig kafkaClientConfig) { // Thêm KafkaClientConfig
        this.messagingTemplate = messagingTemplate;
        this.objectMapper = objectMapper;
        this.tripService = tripService;
        // Tạo consumer thủ công từ config đã hoạt động
        this.consumer = kafkaClientConfig.createKafkaConsumer(GROUP_ID);
    }
    
    // === BẮT ĐẦU: Thêm các phương thức quản lý vòng đời ===
    @PostConstruct
    public void start() {
        running.set(true);
        executorService.submit(this::pollLoop);
    }

    public void pollLoop() {
        consumer.subscribe(Collections.singletonList(TOPIC));
        logger.info("TripStatusListener started polling topic '{}'...", TOPIC);
        
        while (running.get()) {
            try {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));
                for (ConsumerRecord<String, String> record : records) {
                    processMessage(record.value());
                }
                if (!records.isEmpty()) {
                    consumer.commitSync();
                }
            } catch (Exception e) {
                logger.error("Error in TripStatusListener polling loop", e);
            }
        }
    }

    @PreDestroy
    public void shutdown() {
        running.set(false);
        executorService.shutdown();
        consumer.close(Duration.ofSeconds(5));
        logger.info("TripStatusListener shut down.");
    }
    // === KẾT THÚC ===
    
    // Phương thức xử lý logic, được gọi từ pollLoop
    public void processMessage(String message) {
        try {
            logger.info("Received message from Kafka topic '{}': {}", TOPIC, message);
            TripStatusUpdateEvent event = objectMapper.readValue(message, TripStatusUpdateEvent.class);

            if (event.getOldStatus() != null && event.getOldStatus() == 0 && event.getNewStatus() == 1) {
                logger.info("Condition met: Trip {} has departed (0 -> 1). Fetching updated trip list.", event.getTripId());

                List<Trip> pendingTrips = tripService.getTripsByStatus(0);
                List<Trip> departedTrips = tripService.getTripsByStatus(1);
                
                List<Trip> tripsToSend = new ArrayList<>();
                tripsToSend.addAll(pendingTrips);
                tripsToSend.addAll(departedTrips);

                messagingTemplate.convertAndSend("/topic/trip-updates", tripsToSend);
                
                logger.info("Sent updated trip list to WebSocket topic '/topic/trip-updates'");
            }
        } catch (Exception e) {
            logger.error("Could not deserialize or process Kafka message: " + e.getMessage());
        }
    }
}