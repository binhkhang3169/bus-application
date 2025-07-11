package com.example.trip_service.controller;

import com.example.trip_service.model.Route;
// import com.example.trip_service.model.SpecialDay; // Not used in this controller
import com.example.trip_service.model.log.Log_route;
import com.example.trip_service.repository.log.LogRouteRepository;
// import com.example.trip_service.repository.log.LogSpecialDayRepository; // Not used in this controller
import com.example.trip_service.service.intef.RouteService;
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
@RequestMapping("/api/v1/routes") // Base path matches service registration in Gateway
public class RouteController {

    private static final Logger log = LoggerFactory.getLogger(RouteController.class);
    private final RouteService service;
    private final LogRouteRepository logRouteRepository; // Make final and use constructor injection

    @Autowired
    public RouteController(RouteService service, LogRouteRepository logRouteRepository) {
        this.service = service;
        this.logRouteRepository = logRouteRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Route> routes = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Routes retrieved successfully",
                "data", routes
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Route r = service.findById(id);
        if (r != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Route found",
                    "data", r
            ));
        } else {
            log.warn("Route not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Route not found",
                    "data", ""
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Route route,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create route request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                route.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating route: {}", userIdString, e.getMessage());
                // Optionally, return a 400 Bad Request if createdBy is mandatory
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating route. 'createdBy' will be null.");
        }
        
        route.setCreatedAt(new Date()); // Set creation timestamp server-side

        Route createdRoute = service.save(route);
        // Logging of creation can be done here if necessary, or rely on entity's createdBy/At fields
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Route created successfully",
                "data", createdRoute
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody Route route) {

        log.info("Received update route request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        Route existingRoute = service.findById(id);
        if (existingRoute == null) {
            log.warn("Route not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Route not found for update",
                    "data", ""
            ));
        }

        route.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If your Route entity has an 'updatedBy' field, set it:
                // route.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating route: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating route id: {}. 'updatedBy' for log will be null.", id);
        }

        // Preserve original creation details if not meant to be updated
        // route.setCreatedAt(existingRoute.getCreatedAt());
        // route.setCreatedBy(existingRoute.getCreatedBy());
        // If Route entity has an 'updatedAt' field:
        // route.setUpdatedAt(new Date());

        Route updatedRoute = service.save(route);

        if (updatedBy != null) { // Log only if we have the user ID
            Log_route logEntry = new Log_route();
            logEntry.setRouteId(id);
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logRouteRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Route updated successfully",
                "data", updatedRoute
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete route request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        
        // Primary authorization should be at the Gateway.
        // If additional checks are needed here based on role for deletion:
        // if (userRole == null || !userRole.equals("ROLE_ADMIN")) { // Example check
        //    log.warn("User with role {} (ID: {}) attempted to delete route {}, forbidden.", userRole, userIdString, id);
        //    return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of("code", HttpStatus.FORBIDDEN.value(), "message", "Insufficient permissions to delete route"));
        // }

        service.deleteById(id); // Add checks in service or here if it doesn't exist
        
        // Log deletion if userId is available
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_route logEntry = new Log_route();
                logEntry.setRouteId(id);
                // logEntry.setAction("DELETED"); // If your log model supports this
                logEntry.setUpdatedAt(new Date()); // Signifying time of deletion
                logEntry.setUpdatedBy(deletedBy);  // Signifying who deleted
                logRouteRepository.save(logEntry);
                log.info("Route id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for route id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Route deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {

        log.info("Request for routes by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        
        // Gateway handles primary RBAC. If request reaches here, user is permitted.
        // Additional logging or fine-grained logic based on role can be added if needed.
        // For example, log which admin/operator accessed it:
        // if (userRole != null && (userRole.equals("ROLE_ADMIN") || userRole.equals("ROLE_OPERATOR"))) {
        //     log.info("Privileged UserID: {} with Role: {} is accessing routes by status.", userIdString, userRole);
        // }

        List<Route> routes = service.getRoutesByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Routes retrieved successfully by status",
                "data", routes
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for route id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of route id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating status of route id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Ensure service.updateRouteStatus can accept actorUserId (updatedBy)
            // Original: service.updateRouteStatus(userId,id, status) where userId was @RequestParam
            // New: service.updateRouteStatus(actorUserId, routeId, newStatus)
            boolean result = service.updateRouteStatus(updatedBy, id, status);

            if (result) {
                // Log the status update if needed, potentially with old/new status
                if (updatedBy != null) {
                    Log_route logEntry = new Log_route();
                    logEntry.setRouteId(id);
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    // logEntry.setOldStatus(...); logEntry.setNewStatus(status); // If desired
                    logRouteRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Route status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Route not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Route not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or other service-specific exceptions
            log.error("Error updating route status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating route status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
}