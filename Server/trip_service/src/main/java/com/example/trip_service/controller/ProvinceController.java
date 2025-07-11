package com.example.trip_service.controller;

import com.example.trip_service.model.Province;
// import com.example.trip_service.model.Route; // Not used in this controller
import com.example.trip_service.model.log.Log_province;
import com.example.trip_service.repository.log.LogProvinceRepository;
// import com.example.trip_service.repository.log.LogRouteRepository; // Not used in this controller
// import com.example.trip_service.security.JwtTokenUtil; // To be removed or adapted if still needed for other purposes
import com.example.trip_service.service.intef.ProvinceService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Date;
// import java.util.HashMap; // Not used
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/provinces") // Base path matches service registration in Gateway
public class ProvinceController {

    private static final Logger log = LoggerFactory.getLogger(ProvinceController.class);
    private final ProvinceService service;
    private final LogProvinceRepository logProvinceRepository;

    @Autowired
    public ProvinceController(ProvinceService service, LogProvinceRepository logProvinceRepository) {
        this.service = service;
        this.logProvinceRepository = logProvinceRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Province> list = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Provinces retrieved successfully",
                "data", list
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Province p = service.findById(id);
        if (p != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Province found",
                    "data", p
            ));
        } else {
            log.warn("Province not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Province not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Province province,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) { // Optional: get role if needed for logic

        log.info("Received create province request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                province.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating province: {}", userIdString, e.getMessage());
                // Decide if this is a bad request or if createdBy can be null
                // return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("code", HttpStatus.BAD_REQUEST.value(), "message", "Invalid X-User-ID format"));
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating province. 'createdBy' will be null.");
            // If createdBy is mandatory, return an error here.
        }
        
        province.setCreatedAt(new Date()); // Set creation timestamp server-side

        Province saved = service.save(province);
        // Log creation if needed, though 'createdBy' is now on the entity
        // Log_province logEntry = new Log_province();
        // logEntry.setProvinceId(saved.getId());
        // logEntry.setCreatedAt(saved.getCreatedAt()); // or new Date()
        // logEntry.setCreatedBy(createdBy);
        // logProvinceRepository.save(logEntry); // If you have a separate log for creates

        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Province created successfully",
                "data", saved
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole, // Optional
            @PathVariable Integer id,
            @RequestBody Province province) {

        log.info("Received update province request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        // Optional: Check if province exists first (service layer should ideally handle this)
        Province existingProvince = service.findById(id);
        if (existingProvince == null) {
            log.warn("Province not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Province not found for update",
                    "data", ""
            ));
        }

        province.setId(id); // Ensure ID from path is used for the update

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If Province entity has an 'updatedBy' field, set it here:
                // province.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for updating province: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating province id: {}. 'updatedBy' for log will be null.", id);
        }

        // Preserve original createdAt and createdBy unless explicitly allowed to change
        // province.setCreatedAt(existingProvince.getCreatedAt());
        // province.setCreatedBy(existingProvince.getCreatedBy());
        // province.setUpdatedAt(new Date()); // If Province entity has an 'updatedAt' field

        Province updated = service.save(province); // Assuming save handles create or update

        if (updatedBy != null) { // Log only if we have the user ID
            Log_province logEntry = new Log_province();
            logEntry.setProvinceId(id);
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logProvinceRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Province updated successfully",
                "data", updated
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // For logging
            @RequestHeader(name = "X-User-Role", required = false) String userRole, // Optional
            @PathVariable Integer id) {
        
        log.info("Received delete province request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        // Optional: Add logic to check if the role has permission to delete,
        // though this should primarily be handled by the Gateway's RBAC.
        // if (userRole == null || !userRole.equals("ROLE_ADMIN")) {
        //    log.warn("User with role {} attempted to delete province {}, forbidden.", userRole, id);
        //    return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions"));
        // }

        service.deleteById(id); // Consider if this should return a boolean or throw if not found
        
        // Log deletion
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_province logEntry = new Log_province(); // Or a generic action log
                logEntry.setProvinceId(id);
                // You might want a different log structure or an "action" field
                logEntry.setUpdatedAt(new Date()); // Using UpdatedAt to signify the time of action
                logEntry.setUpdatedBy(deletedBy);  // Signifying who performed the action
                // logEntry.setAction("DELETED"); // If your log model supports this
                logProvinceRepository.save(logEntry);
                log.info("Province id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for province id: {}: {}", id, e.getMessage());
            }
        }


        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Province deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            // @RequestHeader("Authorization") String token, // REMOVE THIS - Gateway handles auth
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // For logging or if needed
            @RequestHeader(name = "X-User-Role", required = false) String userRole,   // Role from Gateway
            @PathVariable Integer status) {

        log.info("Request for provinces by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);

        // The API Gateway should have already performed RBAC.
        // If this endpoint is configured in the Gateway to only allow specific roles (e.g., ROLE_ADMIN, ROLE_OPERATOR),
        // then we can assume the user has the necessary role if the request reaches here.
        // No need to re-parse JWT and check role here unless there's very nuanced sub-logic.

        // Example: If for some reason, even after gateway RBAC, you need to log which admin accessed it:
        // if (userRole != null && (userRole.equals("ROLE_ADMIN") || userRole.equals("ROLE_OPERATOR"))) {
        //     log.info("Admin/Operator UserID: {} with Role: {} is accessing provinces by status.", userIdString, userRole);
        // } else {
        //     // This case should ideally not happen if gateway RBAC is correct.
        //     log.warn("User with unexpected role {} accessed /status/{status}. Gateway RBAC might need review.", userRole);
        //     // return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions"));
        // }


        List<Province> provinces = service.getProvincesByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Provinces retrieved successfully by status",
                "data", provinces
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // For updatedBy in log
            @RequestHeader(name = "X-User-Role", required = false) String userRole, // Optional
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for province id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of province id: {}: {}",userIdString, id, e.getMessage());
            }
        } else {
             log.warn("X-User-ID header not present or empty for updating status of province id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Assuming service.updateProvinceStatus now correctly takes the actor's ID for its internal logic/logging.
            // The original method took (@RequestParam Integer userId, ...).
            // It should be changed to (Integer actorUserId, Integer provinceId, Integer newStatus)
            boolean result = service.updateProvinceStatus(updatedBy, id, status);

            if (result) {
                // Log separately here if the service method doesn't log with updatedBy,
                // or if you want a controller-level log.
                // The existing service method `updateProvinceStatus(userId,id, status)` suggests it might already log or use the userId.
                // If it does, ensure it uses the `updatedBy` passed to it now.
                // If it does NOT, and you need to log it using LogProvinceRepository:
                if (updatedBy != null) {
                    Log_province logEntry = new Log_province();
                    logEntry.setProvinceId(id);
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    // logEntry.setOldStatus(existingProvince.getStatus()); // If you fetch existing first
                    // logEntry.setNewStatus(status);
                    logProvinceRepository.save(logEntry);
                }

                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Province status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Province not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Province not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or any other specific exception your service might throw
            log.error("Error updating province status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of( // Or NOT_FOUND if appropriate
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating province status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
}