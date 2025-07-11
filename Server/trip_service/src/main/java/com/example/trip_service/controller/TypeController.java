package com.example.trip_service.controller;

import com.example.trip_service.model.Type;
// import com.example.trip_service.model.Vehicle; // Not used in this controller
import com.example.trip_service.model.log.Log_type;
import com.example.trip_service.repository.log.LogTypeRepository;
import com.example.trip_service.service.intef.TypeService;
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
@RequestMapping("/api/v1/types") // Base path matches service registration in Gateway
public class TypeController {

    private static final Logger log = LoggerFactory.getLogger(TypeController.class);
    private final TypeService service;
    private final LogTypeRepository logTypeRepository; // Make final

    @Autowired
    public TypeController(TypeService service, LogTypeRepository logTypeRepository) {
        this.service = service;
        this.logTypeRepository = logTypeRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Type> types = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Types retrieved successfully",
                "data", types
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Type t = service.findById(id);
        if (t != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Type found",
                    "data", t
            ));
        } else {
            log.warn("Type not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Type not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Type type, // Renamed parameter from your previous example
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create type request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                type.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating type: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating type. 'createdBy' will be null.");
        }
        
        type.setCreatedAt(new Date()); // Set creation timestamp server-side

        Type createdType = service.save(type);
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Type created successfully",
                "data", createdType
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody Type type) { // Renamed parameter

        log.info("Received update type request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        Type existingType = service.findById(id);
        if (existingType == null) {
            log.warn("Type not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Type not found for update",
                    "data", ""
            ));
        }

        type.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If Type entity has 'updatedBy' field:
                // type.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating type: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating type id: {}. 'updatedBy' for log will be null.", id);
        }
        
        // Preserve original creation details if not meant to be updated
        // type.setCreatedAt(existingType.getCreatedAt());
        // type.setCreatedBy(existingType.getCreatedBy());
        // If Type entity has 'updatedAt' field:
        // type.setUpdatedAt(new Date());

        Type updatedType = service.save(type);

        if (updatedBy != null) {
            Log_type logEntry = new Log_type(); // Ensure this is your correct log model
            logEntry.setTypeId(id); // Assuming this setter exists
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logTypeRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Type updated successfully",
                "data", updatedType
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete type request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        
        // Primary authorization is at the Gateway.
        // If additional checks based on role are needed for deletion:
        // if (userRole == null || !userRole.equals("ROLE_ADMIN")) { // Example
        //    log.warn("User with role {} attempted to delete type {}, forbidden.", userRole, id);
        //    return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions"));
        // }

        service.deleteById(id); // Ensure service layer handles "not found" or check here
        
        // Log deletion
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_type logEntry = new Log_type();
                logEntry.setTypeId(id); // Assuming setter exists
                logEntry.setUpdatedAt(new Date()); // Or a 'deletedAt' field
                logEntry.setUpdatedBy(deletedBy);  // Or 'deletedBy'
                logTypeRepository.save(logEntry);
                log.info("Type id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for type id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Type deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {

        log.info("Request for types by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        
        // Gateway handles primary RBAC.

        List<Type> types = service.getTypesByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Types retrieved successfully by status",
                "data", types
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for type id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of type id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
             log.warn("X-User-ID header not present or empty for updating status of type id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Ensure service.updateTypeStatus can accept actorUserId (updatedBy)
            // Original: service.updateTypeStatus(userId, id, status) where userId was @RequestParam
            boolean result = service.updateTypeStatus(updatedBy, id, status);

            if (result) {
                if (updatedBy != null) {
                    Log_type logEntry = new Log_type();
                    logEntry.setTypeId(id); // Assuming setter exists
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    logTypeRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Type status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Type not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Type not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or other service-specific exceptions
            log.error("Error updating type status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating type status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
}