package com.example.trip_service.service.impl;

import java.time.Duration;
import java.util.Arrays;
import java.util.Optional;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicBoolean;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import com.example.trip_service.config.KafkaClientConfig;
import com.example.trip_service.dto.SeatUpdateEvent;
import com.example.trip_service.model.Trip;
import com.example.trip_service.repository.TripRepository;
import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Component
public class SeatUpdateListener {

    private static final Logger logger = LoggerFactory.getLogger(SeatUpdateListener.class);
    private final ObjectMapper objectMapper;
    private final TripRepository tripRepository;
    private final KafkaConsumer<String, String> consumer;
    private final ExecutorService executorService = Executors.newSingleThreadExecutor();
    private final AtomicBoolean running = new AtomicBoolean(false);

    private static final String SEATS_RESERVED_TOPIC = "seats_reserved";
    private static final String SEATS_RELEASED_TOPIC = "seats_released";
    private static final String GROUP_ID = "trip_seat_updater";

    @Autowired
    public SeatUpdateListener(ObjectMapper objectMapper, TripRepository tripRepository,
            KafkaClientConfig kafkaClientConfig) {
        this.objectMapper = objectMapper;
        this.tripRepository = tripRepository;
        this.consumer = kafkaClientConfig.createKafkaConsumer(GROUP_ID);
    }

    @PostConstruct
    public void start() {
        running.set(true);
        executorService.submit(this::pollLoop);
    }

    public void pollLoop() {
        consumer.subscribe(Arrays.asList(SEATS_RESERVED_TOPIC, SEATS_RELEASED_TOPIC));
        logger.info("SeatUpdateListener started polling...");
        while (running.get()) {
            try {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(1000));
                for (ConsumerRecord<String, String> record : records) {
                    logger.info("Received record from topic {}: {}", record.topic(), record.value());
                    if (record.topic().equals(SEATS_RESERVED_TOPIC)) {
                        handleSeatsReserved(record.value());
                    } else if (record.topic().equals(SEATS_RELEASED_TOPIC)) {
                        handleSeatsReleased(record.value());
                    }
                }
                if (!records.isEmpty()) {
                    consumer.commitSync(); // Commit offset sau khi xử lý xong batch
                }
            } catch (Exception e) {
                logger.error("Error in Kafka polling loop", e);
            }
        }
    }

    // Logic xử lý vẫn giữ nguyên, nhưng giờ được gọi từ pollLoop
    @Transactional
    public void handleSeatsReserved(String message) {
        try {
            SeatUpdateEvent event = objectMapper.readValue(message, SeatUpdateEvent.class);
            logger.info("Processing SEATS_RESERVED for TripID: {}", event.getTripId());
            Optional<Trip> tripOptional = tripRepository.findById(Integer.parseInt(event.getTripId()));
            if (tripOptional.isPresent()) {
                Trip trip = tripOptional.get();
                int currentStock = trip.getStock();
                trip.setStock(currentStock - event.getSeatCount());
                tripRepository.save(trip);
                logger.info("Updated stock for TripID {}. New stock: {}", trip.getId(), trip.getStock());
            } else {
                logger.warn("TripID {} not found for seat reservation.", event.getTripId());
            }
        } catch (Exception e) {
            logger.error("Error processing SEATS_RESERVED message: " + message, e);
        }
    }

    @Transactional
    public void handleSeatsReleased(String message) {
        try {
            SeatUpdateEvent event = objectMapper.readValue(message, SeatUpdateEvent.class);
            logger.info("Received SEATS_RELEASED event for TripID: {}, SeatCount: {}", event.getTripId(),
                    event.getSeatCount());

            Optional<Trip> tripOptional = tripRepository.findById(Integer.parseInt(event.getTripId()));

            if (tripOptional.isPresent()) {
                Trip trip = tripOptional.get();
                trip.setStock(trip.getStock() + event.getSeatCount()); // Increment the stock

                tripRepository.save(trip);
                logger.info("Successfully released seats for TripID {}. New stock: {}", trip.getId(), trip.getStock());
            } else {
                logger.warn("Received seat release for non-existent TripID: {}", event.getTripId());
            }
        } catch (Exception e) {
            logger.error("Error processing SEATS_RELEASED message: " + message, e);
        }
    }

    @PreDestroy
    public void shutdown() {
        running.set(false);
        executorService.shutdown();
        consumer.close(Duration.ofSeconds(5));
        logger.info("SeatUpdateListener shut down.");
    }
}