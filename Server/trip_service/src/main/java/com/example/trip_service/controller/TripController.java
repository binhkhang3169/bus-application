package com.example.trip_service.controller;

import java.util.Date;
import java.util.List;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.example.trip_service.dto.TripInfoProjection;
import com.example.trip_service.model.Trip;
import com.example.trip_service.model.log.Log_trip;
import com.example.trip_service.repository.log.LogTripRepository;
import com.example.trip_service.service.intef.TripService;

@RestController
@RequestMapping("/api/v1/trips")
public class TripController {

    private static final Logger log = LoggerFactory.getLogger(TripController.class);
    private final TripService service;
    private final LogTripRepository logTripRepository; // Make final

    @Autowired
    public TripController(TripService service, LogTripRepository logTripRepository) {
        this.service = service;
        this.logTripRepository = logTripRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Trip> trips = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trips retrieved successfully",
                "data", trips
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Trip t = service.findById(id);
        if (t != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Trip found",
                    "data", t
            ));
        } else {
            log.warn("Trip not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Trip not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Trip trip,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create trip request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                trip.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating trip: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating trip. 'createdBy' will be null.");
        }

        trip.setCreatedAt(new Date()); // Set creation timestamp server-side

        Trip createdTrip = service.save(trip);
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Trip created successfully",
                "data", createdTrip
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody Trip trip) {

        log.info("Received update trip request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        Trip existingTrip = service.findById(id);
        if (existingTrip == null) {
            log.warn("Trip not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Trip not found for update",
                    "data", ""
            ));
        }

        trip.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If Trip entity has 'updatedBy' field:
                // trip.setUpdatedBy(updatedBy); // Assuming Trip model has setUpdatedBy
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating trip: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating trip id: {}. 'updatedBy' for log will be null.", id);
        }

        // Preserve original creation details if not meant to be updated by request body
        // trip.setCreatedAt(existingTrip.getCreatedAt());
        // trip.setCreatedBy(existingTrip.getCreatedBy());
        // If Trip entity has 'updatedAt' field:
        // trip.setUpdatedAt(new Date()); // Assuming Trip model has setUpdatedAt

        Trip updatedTrip = service.save(trip);

        if (updatedBy != null) {
            Log_trip logEntry = new Log_trip();
            logEntry.setTripId(id);
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logTripRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trip updated successfully",
                "data", updatedTrip
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete trip request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        service.deleteById(id); // Ensure service layer handles "not found" or check here

        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_trip logEntry = new Log_trip();
                logEntry.setTripId(id);
                logEntry.setUpdatedAt(new Date()); // Or a specific 'deletedAt' field
                logEntry.setUpdatedBy(deletedBy);  // Or 'deletedBy'
                logTripRepository.save(logEntry);
                log.info("Trip id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for trip id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trip deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/search")
    public ResponseEntity<Map<String, Object>> searchTrips(
            @RequestParam String from,
            @RequestParam Integer fromId,
            @RequestParam(required = false) String fromTime,
            @RequestParam String to,
            @RequestParam Integer toId,
            @RequestParam(required = false) String toTime,
            @RequestParam(required = false, defaultValue = "1") Integer quantity,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString
    ) {
        Integer userId = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                userId = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for trip search: {}", userIdString, e.getMessage());
            }
        }

        log.info("Search trips request: fromId={}, toId={}, fromTime={}, toTime={}, quantity={}, userId={}", 
                 fromId, toId, fromTime, toTime, quantity, userId);
        
        String departureDate = fromTime != null ? fromTime : toTime;
        
        List<TripInfoProjection> trips = service.searchTripsByLocationsWithSeats(fromId, toId, departureDate, quantity, userId);

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trip search successful",
                "data", trips
        ));
    }

    @GetMapping("/{id}/seats")
    public ResponseEntity<Map<String, Object>> getTripWithSeats(@PathVariable int id) {
        log.info("Request for trip with seats, id: {}", id);
        TripInfoProjection trip = service.findByIdWithSeats(id);
        if (trip == null) {
            log.warn("Trip with seats not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Trip with seats not found",
                    "data", ""
            ));
        }
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trip with seats found",
                "data", trip
        ));
    }

    @GetMapping("/search/basic")
    public ResponseEntity<Map<String, Object>> searchTripsBasic(
            @RequestParam String from,
            @RequestParam Integer fromId,
            @RequestParam(required = false) String fromTime,
            @RequestParam String to,
            @RequestParam Integer toId,
            @RequestParam(required = false) String toTime
    ) {
        log.info("Basic search trips request: fromId={}, toId={}, fromTime={}, toTime={}", fromId, toId, fromTime, toTime);
        String departureDate = fromTime != null ? fromTime : toTime;
        Integer quantity = 1; // Default quantity

        List<TripInfoProjection> trips = service.searchTripsByLocations(fromId, toId, departureDate, quantity);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Basic trip search successful",
                "data", trips
        ));
    }


    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {

        log.info("Request for trips by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        List<Trip> trips = service.getTripsByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Trips retrieved successfully by status",
                "data", trips
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for trip id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of trip id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
             log.warn("X-User-ID header not present or empty for updating status of trip id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            boolean result = service.updateTripStatus(updatedBy, id, status);

            if (result) {
                if (updatedBy != null) { // Log only if user ID was successfully parsed
                    Log_trip logEntry = new Log_trip();
                    logEntry.setTripId(id);
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    logTripRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Trip status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Trip not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Trip not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) {
            log.error("Error updating trip status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating trip status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }

    // Modified endpoint: Get trips for the driver identified by X-User-ID
    @GetMapping("/driver")
    public ResponseEntity<Map<String, Object>> getAssignedTripsForDriver(
            @RequestHeader(name = "X-User-ID") String userIdString, // Made required for this endpoint
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Request for assigned trips. UserID (as DriverID): {}, Role: {}", userIdString, userRole);

        Integer driverId;
        try {
            driverId = Integer.parseInt(userIdString);
        } catch (NumberFormatException e) {
            log.warn("Invalid X-User-ID format: {}. Expected an integer.", userIdString);
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Invalid User ID format in X-User-ID header. Expected an integer.",
                    "data", ""
            ));
        }

        List<Trip> trips = service.findTripsByDriverId(driverId);

        if (trips.isEmpty()) {
            log.info("No assigned trips found for driverId: {}", driverId);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", trips.isEmpty() ? "No assigned trips found" : "Assigned trips retrieved successfully",
                "data", trips
        ));
    }
}