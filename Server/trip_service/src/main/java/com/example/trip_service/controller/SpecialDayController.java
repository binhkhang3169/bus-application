package com.example.trip_service.controller;

import com.example.trip_service.model.SpecialDay;
// import com.example.trip_service.model.Station; // Not used in this controller
import com.example.trip_service.model.log.Log_specialDay; // Assuming this is the correct log model name
import com.example.trip_service.repository.log.LogSpecialDayRepository;
import com.example.trip_service.service.intef.SpecialDayService;
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
@RequestMapping("/api/v1/special-days") // Base path matches service registration in Gateway
public class SpecialDayController {

    private static final Logger log = LoggerFactory.getLogger(SpecialDayController.class);
    private final SpecialDayService service;
    private final LogSpecialDayRepository logSpecialDayRepository; // Make final

    @Autowired
    public SpecialDayController(SpecialDayService service, LogSpecialDayRepository logSpecialDayRepository) {
        this.service = service;
        this.logSpecialDayRepository = logSpecialDayRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<SpecialDay> days = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Special days retrieved successfully",
                "data", days
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        SpecialDay d = service.findById(id);
        if (d != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Special day found",
                    "data", d
            ));
        } else {
            log.warn("SpecialDay not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Special day not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody SpecialDay day,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create special day request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                day.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating special day: {}", userIdString, e.getMessage());
                // Optionally, return a 400 Bad Request if createdBy is mandatory
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating special day. 'createdBy' will be null.");
        }
        
        day.setCreatedAt(new Date()); // Set creation timestamp server-side

        SpecialDay createdDay = service.save(day);
        // Logging of creation can be done here or rely on entity's createdBy/At
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Special day created successfully",
                "data", createdDay
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody SpecialDay day) {

        log.info("Received update special day request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        SpecialDay existingDay = service.findById(id);
        if (existingDay == null) {
            log.warn("SpecialDay not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Special day not found for update",
                    "data", ""
            ));
        }

        day.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If SpecialDay entity has 'updatedBy' field:
                // day.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating special day: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating special day id: {}. 'updatedBy' for log will be null.", id);
        }
        
        // Preserve original creation details if not meant to be updated
        // day.setCreatedAt(existingDay.getCreatedAt());
        // day.setCreatedBy(existingDay.getCreatedBy());
        // If SpecialDay entity has 'updatedAt' field:
        // day.setUpdatedAt(new Date());

        SpecialDay updatedDay = service.save(day);

        if (updatedBy != null) { // Log only if we have the user ID
            Log_specialDay logEntry = new Log_specialDay(); // Ensure correct log model name
            logEntry.setSpecialDayId(id); // Assuming this setter exists
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logSpecialDayRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Special day updated successfully",
                "data", updatedDay
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete special day request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        
        // Primary authorization is at the Gateway.
        // Additional role check here if needed:
        // if (userRole == null || !userRole.equals("ROLE_ADMIN")) { // Example
        //    log.warn("User with role {} attempted to delete special day {}, forbidden.", userRole, id);
        //    return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions"));
        // }

        service.deleteById(id); // Ensure service layer handles "not found" appropriately or check here
        
        // Log deletion
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_specialDay logEntry = new Log_specialDay();
                logEntry.setSpecialDayId(id); // Assuming this setter exists
                // logEntry.setAction("DELETED"); // If your log model supports this
                logEntry.setUpdatedAt(new Date()); // Using UpdatedAt to signify time of deletion
                logEntry.setUpdatedBy(deletedBy);  // Signifying who performed the action
                logSpecialDayRepository.save(logEntry);
                log.info("Special day id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for special day id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Special day deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {

        log.info("Request for special days by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        
        // Gateway handles primary RBAC.
        // If this endpoint (e.g., /api/v1/special-days/status/1) is configured in the Gateway
        // to require ROLE_ADMIN, then requests reaching here have already passed that check.

        List<SpecialDay> specialDays = service.getSpecialDaysByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Special days retrieved successfully by status",
                "data", specialDays
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for special day id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of special day id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating status of special day id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Ensure service.updateSpecialDayStatus can accept actorUserId (updatedBy)
            // Original: service.updateSpecialDayStatus(userId, id, status) where userId was @RequestParam
            // New: service.updateSpecialDayStatus(actorUserId, specialDayId, newStatus)
            boolean result = service.updateSpecialDayStatus(updatedBy, id, status);

            if (result) {
                if (updatedBy != null) {
                    Log_specialDay logEntry = new Log_specialDay();
                    logEntry.setSpecialDayId(id); // Assuming this setter exists
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    logSpecialDayRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Special day status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Special day not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Special day not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or other service-specific exceptions
            log.error("Error updating special day status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating special day status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
}