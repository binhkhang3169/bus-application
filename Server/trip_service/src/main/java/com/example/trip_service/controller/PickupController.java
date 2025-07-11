package com.example.trip_service.controller;

import com.example.trip_service.model.Pickup;
// import com.example.trip_service.model.Province; // Not used in this controller directly
import com.example.trip_service.model.log.Log_pickup;
import com.example.trip_service.repository.log.LogPickupRepository;
// import com.example.trip_service.repository.log.LogProvinceRepository; // Not used in this controller directly
import com.example.trip_service.service.intef.PickupService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Date;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/pickups") // This base path is what the API Gateway will proxy to
public class PickupController {

    private final PickupService service;
    private final LogPickupRepository logPickupRepository; // Make final and use constructor injection

    @Autowired
    public PickupController(PickupService service, LogPickupRepository logPickupRepository) {
        this.service = service;
        this.logPickupRepository = logPickupRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Pickup> pickups = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Pickups retrieved successfully",
                "data", pickups
        ));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable String id) {
        Pickup p = service.findById(id);
        if (p != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Pickup found",
                    "data", p
            ));
        } else {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Pickup not found",
                    "data", "" // Consider returning null or an empty object for data
            ));
        }
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Pickup pickup,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString) { // Read X-User-ID header

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                pickup.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                // Log the error, and decide if this should be a bad request
                System.err.println("Warning: Could not parse X-User-ID header value: " + userIdString + " - " + e.getMessage());
                // Optionally, return a 400 Bad Request if createdBy is mandatory and parsing fails
                // return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("code", HttpStatus.BAD_REQUEST.value(), "message", "Invalid X-User-ID format"));
            }
        } else {
            // Handle case where X-User-ID is not present, if necessary
            // For example, if it's optional or for system-generated entities
             System.err.println("Warning: X-User-ID header not present or empty for POST /api/v1/pickups");
        }
        
        // It's good practice to set createdAt on the server-side
        pickup.setCreatedAt(new Date()); 

        Pickup createdPickup = service.save(pickup);
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Pickup created successfully",
                "data", createdPickup
        ));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // Read X-User-ID for updatedBy
            @PathVariable String id,
            @RequestBody Pickup pickup) {

        // First, check if the pickup exists
        Pickup existingPickup = service.findById(id);
        if (existingPickup == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Pickup not found for update with id: " + id,
                    "data", ""
            ));
        }
        
        pickup.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // You might want to set an 'updatedBy' field on the Pickup entity itself if it exists
                // pickup.setUpdatedBy(updatedBy); 
            } catch (NumberFormatException e) {
                System.err.println("Warning: Could not parse X-User-ID header value for update: " + userIdString + " - " + e.getMessage());
            }
        } else {
             System.err.println("Warning: X-User-ID header not present or empty for PUT /api/v1/pickups/" + id);
        }

        // Retain original createdBy and createdAt unless they are part of the update payload
        // and you intend to allow their modification (usually not for createdBy/At).
        // If Pickup entity has its own createdBy, it might be good to preserve it:
        // pickup.setCreatedBy(existingPickup.getCreatedBy());
        // pickup.setCreatedAt(existingPickup.getCreatedAt());


        Pickup updatedPickup = service.save(pickup); // The save method should handle update logic

        // Logging the update action
        if (updatedBy != null) { // Only log if updatedBy is available
            Log_pickup logPickup = new Log_pickup();
            logPickup.setPickupId(id);
            logPickup.setUpdatedAt(new Date());
            logPickup.setUpdatedBy(updatedBy); // Use parsed userIdString for log
            logPickupRepository.save(logPickup);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Pickup updated successfully",
                "data", updatedPickup
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // For logging who deleted it
            @PathVariable String id) {
        
        // Optional: Log who performed the delete action
        if (userIdString != null && !userIdString.isEmpty()) {
             try {
                Integer deletedBy = Integer.parseInt(userIdString);
                // Create a log entry for delete if needed
                // Log_pickup logDelete = new Log_pickup();
                // logDelete.setPickupId(id);
                // logDelete.setAction("DELETE"); // You might add an action field to your log
                // logDelete.setUpdatedAt(new Date()); // or a deletedAt field
                // logDelete.setUpdatedBy(deletedBy); // or a deletedBy field
                // logPickupRepository.save(logDelete);
                System.out.println("Info: Pickup " + id + " deleted by UserID: " + deletedBy);
            } catch (NumberFormatException e) {
                System.err.println("Warning: Could not parse X-User-ID for delete log: " + userIdString);
            }
        }

        service.deleteById(id); // Assuming this throws an exception or returns a boolean if not found
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Pickup deleted successfully",
                "data", ""
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(@PathVariable Integer status) {
        List<Pickup> pickups = service.getPickupsByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Pickups retrieved successfully by status", // Message was "Provinces..."
                "data", pickups
        ));
    }

    @PutMapping("/{id}/{status}")
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString, // Read X-User-ID for updatedBy in log
            @PathVariable String id,
            @PathVariable Integer status) {
        
        Integer updatedBy = null;
         if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                 System.err.println("Warning: Could not parse X-User-ID for updateStatus: " + userIdString);
            }
        } else {
            System.err.println("Warning: X-User-ID header not present or empty for PUT /api/v1/pickups/" + id + "/" + status);
        }

        try {
            // Pass updatedBy to the service method if it needs to know who initiated the status update
            // For now, the provided service.updatePickupStatus only takes id and status.
            // If your service layer logs, it might need updatedBy too.
            // The existing `service.updatePickupStatus(userId, id, status)` in your original code
            // seemed to imply `userId` was a @RequestParam. Now it's from the header.
            // We will pass `updatedBy` (parsed from header) to the service if it's designed to accept it.
            // Assuming service.updatePickupStatus now takes (Integer actorUserId, String pickupId, Integer status)
            
            // boolean result = service.updatePickupStatus(updatedBy, id, status); // If service method is changed
            // For now, sticking to the original service method signature shown in your code for updatePickupStatus,
            // but it's odd that it takes a userId as a @RequestParam and also logs internally.
            // Let's assume the service method `updatePickupStatus` already handles logging with the passed userId.
            // The `userId` parameter in the original method signature: `updateStatus(@RequestParam Integer userId, ...)`
            // should now be `updatedBy` (from header).
            // If the service method signature is `updatePickupStatus(String pickupId, Integer newStatus, Integer actorUserId)`
            
            boolean result = service.updatePickupStatus(updatedBy, id, status); // Assuming the service method now accepts the actor's ID this way for its internal logging or operations.

            if (result) {
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Pickup status updated successfully",
                        "data", ""
                ));
            } else {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Pickup not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Assuming service throws this if pickup not found
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Pickup not found: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }
    // New Endpoint
    @GetMapping("/byRoute/{routeId}")
    public ResponseEntity<Map<String, Object>> getPickupsByRouteIdAndSelfIdIsMinusOne(@PathVariable Integer routeId) {
        if (routeId == null) {
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Route ID cannot be null",
                    "data", ""
            ));
        }
        List<Pickup> pickups = service.findByRouteIdAndSelfIdIsMinusOne(routeId);
        if (pickups.isEmpty()) {
            return ResponseEntity.ok(Map.of( // Or NOT_FOUND, depending on preference for empty results
                    "code", HttpStatus.OK.value(),
                    "message", "No pickups found for route ID " + routeId + " with self_id = -1",
                    "data", pickups // will be an empty list
            ));
        }
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Pickups retrieved successfully for route ID " + routeId + " and self_id = -1",
                "data", pickups
        ));
    }
}