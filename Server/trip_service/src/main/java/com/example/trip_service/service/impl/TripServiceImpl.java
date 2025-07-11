package com.example.trip_service.service.impl;

import java.time.LocalDateTime;
import java.util.Date;
import java.util.List;
import java.util.Optional;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.server.ResponseStatusException;

import com.example.trip_service.dto.TripCreatedEvent;
import com.example.trip_service.dto.TripInfoProjection;
import com.example.trip_service.dto.TripSearchEvent;
import com.example.trip_service.dto.TripStatusUpdateEvent;
import com.example.trip_service.model.Trip;
import com.example.trip_service.model.log.Log_trip;
import com.example.trip_service.repository.TripRepository;
import com.example.trip_service.repository.log.LogTripRepository;
import com.example.trip_service.service.intef.TripService;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class TripServiceImpl implements TripService {

    private final TripRepository repository;

    @Autowired
    private LogTripRepository logTripRepository;

    @Autowired
    private RestTemplate restTemplate;

    private final KafkaProducer<String, String> kafkaProducer;
    private final ObjectMapper objectMapper;

    // Tên topic Kafka
    private static final String TRIP_CREATED_TOPIC = "trip_created";
    private static final String TRIP_SEARCH_TOPIC = "trip_search";
    private static final String TRIP_STATUS_UPDATED_TOPIC = "trip_status_updated"; // Topic mới

    @Autowired
    public TripServiceImpl(TripRepository repository,
            KafkaProducer<String, String> kafkaProducer, ObjectMapper objectMapper) {
        this.repository = repository;
        this.kafkaProducer = kafkaProducer;
        this.objectMapper = objectMapper;
    }

    @Override
    public List<Trip> findAll() {
        return repository.findAll();
    }

    @Override
    public Trip findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Trip save(Trip trip) {
        if (trip.getDriverId() != null
                && trip.getDepartureDate() != null && trip.getDepartureTime() != null
                && trip.getArrivalDate() != null && trip.getArrivalTime() != null) {

            LocalDateTime newTripDepartureDateTime = LocalDateTime.of(trip.getDepartureDate(), trip.getDepartureTime());
            LocalDateTime newTripArrivalDateTime = LocalDateTime.of(trip.getArrivalDate(), trip.getArrivalTime());

            LocalDateTime checkRangeStart = newTripDepartureDateTime.minusHours(1);
            LocalDateTime checkRangeEnd = newTripArrivalDateTime.plusHours(1);

            Trip conflictingTrip = repository.findConflictingTripForDriver(
                    trip.getDriverId(),
                    checkRangeStart.toLocalDate(),
                    checkRangeStart.toLocalTime(),
                    checkRangeEnd.toLocalDate(),
                    checkRangeEnd.toLocalTime(),
                    trip.getId() // Pass current trip's ID to exclude it if updating
            );

            if (conflictingTrip != null) {
                throw new ResponseStatusException(HttpStatus.CONFLICT,
                        "Driver has a conflicting trip schedule. Ensure at least a 1-hour break between trips. Conflicting trip ID: "
                                + conflictingTrip.getId());
            }
        } else if (trip.getDriverId() != null) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    "Departure/Arrival date and time are required for driver schedule check.");
        }

        boolean isNewTrip = (trip.getId() == null);

        if (trip.getCreatedAt() == null) {
            trip.setCreatedAt(new Date());
        }

        Trip savedTrip = repository.save(trip);

        if (isNewTrip && savedTrip.getId() != null) {
            sendTripCreationEvent(savedTrip);
        }

        return savedTrip;
    }

    private void sendTripCreationEvent(Trip trip) {
        try {
            TripCreatedEvent event = new TripCreatedEvent(trip.getId().toString(), trip.getTotal(), new Date());
            String eventJson = objectMapper.writeValueAsString(event);
            ProducerRecord<String, String> record = new ProducerRecord<>(TRIP_CREATED_TOPIC, trip.getId().toString(),
                    eventJson);

            kafkaProducer.send(record, (metadata, exception) -> {
                if (exception == null) {
                    System.out.printf("Sent event to topic %s, partition %d, offset %d\n", metadata.topic(),
                            metadata.partition(), metadata.offset());
                } else {
                    System.err.println("Failed to send trip creation event: " + exception.getMessage());
                }
            });
        } catch (Exception e) {
            System.err.println("Error creating or sending trip creation event: " + e.getMessage());
        }
    }

    // Phương thức gửi sự kiện cập nhật trạng thái
    private void sendTripStatusUpdateEvent(Integer tripId, Integer oldStatus, Integer newStatus, Integer updatedBy) {
        try {
            TripStatusUpdateEvent event = new TripStatusUpdateEvent(tripId, oldStatus, newStatus, updatedBy,
                    new Date());
            String eventJson = objectMapper.writeValueAsString(event);
            ProducerRecord<String, String> record = new ProducerRecord<>(TRIP_STATUS_UPDATED_TOPIC, tripId.toString(),
                    eventJson);

            kafkaProducer.send(record, (metadata, exception) -> {
                if (exception == null) {
                    System.out.printf("Sent status update event to topic %s for trip %d\n", metadata.topic(), tripId);
                } else {
                    System.err.println("Failed to send trip status update event: " + exception.getMessage());
                }
            });
        } catch (Exception e) {
            System.err.println("Error creating or sending trip status update event: " + e.getMessage());
        }
    }

    private void sendTripSearchEvent(Integer fromProvinceId, Integer toProvinceId, String departureDate,
            Integer quantity, Integer userId) {
        try {
            // Tạo event với các tham số mới
            TripSearchEvent event = new TripSearchEvent(fromProvinceId, toProvinceId, departureDate, new Date(),
                    quantity, userId);
            String eventJson = objectMapper.writeValueAsString(event);
            ProducerRecord<String, String> record = new ProducerRecord<>(TRIP_SEARCH_TOPIC, eventJson);

            kafkaProducer.send(record, (metadata, exception) -> {
                if (exception == null) {
                    System.out.printf("Sent event to topic %s, partition %d, offset %d\n", metadata.topic(),
                            metadata.partition(), metadata.offset());
                } else {
                    System.err.println("Failed to send trip creation event: " + exception.getMessage());
                }
            });

            // Log chi tiết hơn
            System.out.println("Sent trip search event to Kafka for params: from=" + fromProvinceId + ", to="
                    + toProvinceId + ", date=" + departureDate + ", quantity=" + quantity + ", userId=" + userId);

        } catch (Exception e) {
            System.err.println("Failed to send trip search event. Error: " + e.getMessage());
        }
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }

    @Override
    public List<TripInfoProjection> getTripInfos(Integer routeId, String departureDate, Integer quantity) {
        return repository.findTripsByRouteAndDateRaw(routeId, departureDate, quantity);
    }

    @Override
    public List<TripInfoProjection> searchTripsByLocations(Integer fromProvinceId, Integer toProvinceId,
            String departureDate, Integer quantity) {
        if (departureDate != null && departureDate.contains("-")) {
            String[] parts = departureDate.split("-");
            if (parts.length == 3) {
                departureDate = parts[0] + "-" + parts[1] + "-" + parts[2];
            }
        }
        return repository.findTripsByLocationsAndDate(fromProvinceId, toProvinceId, departureDate, quantity);
    }

    @Override
    public TripInfoProjection findByIdWithSeats(int id) {
        List<TripInfoProjection> tripInfos = repository.findTripInfoById(id);
        if (tripInfos == null || tripInfos.isEmpty()) {
            return null;
        }
        return tripInfos.get(0);
    }

    @Override
    public List<TripInfoProjection> searchTripsByLocationsWithSeats(Integer fromProvinceId, Integer toProvinceId,
            String departureDate, Integer quantity, Integer userId) {

        if (departureDate != null && departureDate.contains("-")) {
            String[] parts = departureDate.split("-");
            if (parts.length == 3) {
                departureDate = parts[0] + "-" + parts[1] + "-" + parts[2];
            }
        }
        // Gọi phương thức đã được cập nhật với các tham số mới
        sendTripSearchEvent(fromProvinceId, toProvinceId, departureDate, quantity, userId);

        return repository.findTripsByLocationsAndDate(fromProvinceId, toProvinceId,
                departureDate, quantity);
    }

    @Override
    public List<Trip> getTripsByStatus(Integer status) {
        // Định nghĩa: 0=Chưa đi, 1=Đã đi, 2=Đã tới, 3=Đã huỷ
        if (status < 0 || status > 3) {
            throw new IllegalArgumentException("Chỉ chấp nhận status từ 0 đến 3");
        }
        return repository.findByStatusIn(List.of(status));
    }

    @Override
    public boolean updateTripStatus(Integer userId, Integer id, Integer newStatus) {
        // Định nghĩa: 0=Chưa đi, 1=Đã đi, 2=Đã tới, 3=Đã huỷ
        if (newStatus < 0 || newStatus > 3) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận giá trị từ 0 đến 3.");
        }
        Optional<Trip> tripOptional = repository.findById(id);
        if (tripOptional.isPresent()) {
            Trip tripUpdate = tripOptional.get();
            Integer oldStatus = tripUpdate.getStatus();

            if (oldStatus.equals(newStatus)) {
                return true; // Không có gì thay đổi
            }

            tripUpdate.setStatus(newStatus);
            repository.save(tripUpdate);

            Log_trip logTrip = new Log_trip();
            logTrip.setTripId(tripUpdate.getId());
            logTrip.setUpdatedAt(new Date());
            logTrip.setUpdatedBy(userId);
            logTripRepository.save(logTrip);

            // Gửi sự kiện tới Kafka
            sendTripStatusUpdateEvent(id, oldStatus, newStatus, userId);

            return true;
        }
        return false;
    }

    @Override
    public List<Trip> findTripsByDriverId(Integer driverId) {
        return repository.findByDriverId(driverId);
    }
}