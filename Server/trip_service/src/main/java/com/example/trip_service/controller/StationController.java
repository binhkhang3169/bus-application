package com.example.trip_service.controller;

import com.example.trip_service.model.Station;
// import com.example.trip_service.model.Trip; // Not used in this controller
import com.example.trip_service.model.log.Log_station;
// import com.example.trip_service.model.log.Log_trip; // Not used in this controller
import com.example.trip_service.repository.log.LogStationRepository;
// import com.example.trip_service.repository.log.LogTripRepository; // Not used in this controller
import com.example.trip_service.service.intef.StationService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Date;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/stations") // Base path matches service registration in Gateway
public class StationController {

    private static final Logger log = LoggerFactory.getLogger(StationController.class);
    private final StationService service;
    private final LogStationRepository logStationRepository; // Make final

    @Autowired
    public StationController(StationService service, LogStationRepository logStationRepository) {
        this.service = service;
        this.logStationRepository = logStationRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Station> stations = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Stations retrieved successfully",
                "data", stations
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Station s = service.findById(id);
        if (s != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Station found",
                    "data", s
            ));
        } else {
            log.warn("Station not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Station not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Station station,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create station request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                station.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating station: {}", userIdString, e.getMessage());
                // Optionally, return a 400 Bad Request if createdBy is mandatory
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating station. 'createdBy' will be null.");
        }
        
        station.setCreatedAt(new Date()); // Set creation timestamp server-side

        Station createdStation = service.save(station);
        // Logging of creation can be done here or rely on entity's createdBy/At fields
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Station created successfully",
                "data", createdStation
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody Station station) {

        log.info("Received update station request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        Station existingStation = service.findById(id);
        if (existingStation == null) {
            log.warn("Station not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Station not found for update",
                    "data", ""
            ));
        }

        station.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If Station entity has 'updatedBy' field:
                // station.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating station: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating station id: {}. 'updatedBy' for log will be null.", id);
        }
        
        // Preserve original creation details if not meant to be updated
        // station.setCreatedAt(existingStation.getCreatedAt());
        // station.setCreatedBy(existingStation.getCreatedBy());
        // If Station entity has 'updatedAt' field:
        // station.setUpdatedAt(new Date());

        Station updatedStation = service.save(station);

        if (updatedBy != null) { // Log only if we have the user ID
            Log_station logEntry = new Log_station(); // Ensure correct log model name
            logEntry.setStationId(id); // Assuming this setter exists
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logStationRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Station updated successfully",
                "data", updatedStation
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete station request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        
        // Primary authorization is at the Gateway.
        // Additional role check here if needed:
        // if (userRole == null || !userRole.equals("ROLE_ADMIN")) { // Example
        //    log.warn("User with role {} attempted to delete station {}, forbidden.", userRole, id);
        //    return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions"));
        // }

        service.deleteById(id); // Ensure service layer handles "not found" appropriately or check here
        
        // Log deletion
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_station logEntry = new Log_station();
                logEntry.setStationId(id); // Assuming this setter exists
                // logEntry.setAction("DELETED"); // If your log model supports this
                logEntry.setUpdatedAt(new Date()); // Using UpdatedAt to signify time of deletion
                logEntry.setUpdatedBy(deletedBy);  // Signifying who performed the action
                logStationRepository.save(logEntry);
                log.info("Station id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for station id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Station deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {

        log.info("Request for stations by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        
        // Gateway handles primary RBAC.
        // If this endpoint is configured in the Gateway to require specific roles,
        // then requests reaching here have already passed that check.

        List<Station> stations = service.getStationsByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Stations retrieved successfully by status",
                "data", stations
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for station id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of station id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating status of station id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Ensure service.updateStationStatus can accept actorUserId (updatedBy)
            // Original: service.updateStationStatus(userId, id, status) where userId was @RequestParam
            // New: service.updateStationStatus(actorUserId, stationId, newStatus)
            boolean result = service.updateStationStatus(updatedBy, id, status);

            if (result) {
                if (updatedBy != null) {
                    Log_station logEntry = new Log_station();
                    logEntry.setStationId(id); // Assuming this setter exists
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    logStationRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Station status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Station not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Station not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or other service-specific exceptions
            log.error("Error updating station status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating station status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
}