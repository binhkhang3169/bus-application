package com.example.trip_service.controller;

import com.example.trip_service.model.Vehicle;
import com.example.trip_service.model.log.Log_vehicle;
import com.example.trip_service.repository.log.LogVehicleRepository;
import com.example.trip_service.service.intef.VehicleService;
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
@RequestMapping("/api/v1/vehicles") // Base path matches service registration in Gateway
public class VehicleController {

    private static final Logger log = LoggerFactory.getLogger(VehicleController.class);
    private final VehicleService service;
    private final LogVehicleRepository logVehicleRepository; // Make final

    @Autowired
    public VehicleController(VehicleService service, LogVehicleRepository logVehicleRepository) {
        this.service = service;
        this.logVehicleRepository = logVehicleRepository;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAll() {
        List<Vehicle> vehicles = service.findAll();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Vehicles retrieved successfully",
                "data", vehicles
        ));
    }

    @GetMapping("/status/{status}")
    public ResponseEntity<Map<String, Object>> getAllByStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer status) {
        
        log.info("Request for vehicles by status: {}. UserID: {}, Role: {}", status, userIdString, userRole);
        // Primary RBAC is at the Gateway.
        List<Vehicle> vehicles = service.getVehiclesByStatus(status);
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Vehicles retrieved successfully by status",
                "data", vehicles
        ));
    }


    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getById(@PathVariable Integer id) {
        Vehicle v = service.findById(id);
        if (v != null) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Vehicle found",
                    "data", v
            ));
        } else {
            log.warn("Vehicle not found for id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Vehicle not found",
                    "data", ""
            ));
        }
    }

    @PutMapping("/{id}/{status}") // For updating status specifically
    public ResponseEntity<Map<String, Object>> updateStatus(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @PathVariable Integer status) {

        log.info("Request to update status for vehicle id: {} to status: {}. UserID: {}, Role: {}", id, status, userIdString, userRole);
        
        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID '{}' for updateStatus of vehicle id: {}: {}", userIdString, id, e.getMessage());
            }
        } else {
             log.warn("X-User-ID header not present or empty for updating status of vehicle id: {}. 'updatedBy' for log/service will be null.", id);
        }

        try {
            // Ensure service.updateVehicleStatus can accept actorUserId (updatedBy)
            // Original: service.updateVehicleStatus(userId, id, status) where userId was @RequestParam
            boolean result = service.updateVehicleStatus(updatedBy, id, status);

            if (result) {
                if (updatedBy != null) {
                    Log_vehicle logEntry = new Log_vehicle();
                    logEntry.setVehicleId(id); // Assuming this setter exists
                    logEntry.setUpdatedAt(new Date());
                    logEntry.setUpdatedBy(updatedBy);
                    logVehicleRepository.save(logEntry);
                }
                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Vehicle status updated successfully",
                        "data", ""
                ));
            } else {
                log.warn("Vehicle not found or status update failed for id: {}", id);
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Vehicle not found or status update failed",
                        "data", ""
                ));
            }
        } catch (IllegalArgumentException ex) { // Or other service-specific exceptions
            log.error("Error updating vehicle status for id: {}: {}", id, ex.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Error updating vehicle status: " + ex.getMessage(),
                    "data", ""
            ));
        }
    }


    @PostMapping
    public ResponseEntity<Map<String, Object>> create(
            @RequestBody Vehicle vehicle,
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole) {

        log.info("Received create vehicle request. UserID: {}, Role: {}", userIdString, userRole);

        Integer createdBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                createdBy = Integer.parseInt(userIdString);
                vehicle.setCreatedBy(createdBy); // Set createdBy from header
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' to Integer for creating vehicle: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for creating vehicle. 'createdBy' will be null.");
        }
        
        vehicle.setCreatedAt(new Date()); // Set creation timestamp server-side

        Vehicle createdVehicle = service.save(vehicle);
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "code", HttpStatus.CREATED.value(),
                "message", "Vehicle created successfully",
                "data", createdVehicle
        ));
    }

    @PutMapping("/{id}") // For updating general vehicle data
    public ResponseEntity<Map<String, Object>> update(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id,
            @RequestBody Vehicle vehicle) {

        log.info("Received update vehicle request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);

        Vehicle existingVehicle = service.findById(id);
        if (existingVehicle == null) {
            log.warn("Vehicle not found for update with id: {}", id);
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Vehicle not found for update",
                    "data", ""
            ));
        }

        vehicle.setId(id); // Ensure ID from path is used

        Integer updatedBy = null;
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                updatedBy = Integer.parseInt(userIdString);
                // If Vehicle entity has 'updatedBy' field:
                // vehicle.setUpdatedBy(updatedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID header '{}' for updating vehicle: {}", userIdString, e.getMessage());
            }
        } else {
            log.warn("X-User-ID header not present or empty for updating vehicle id: {}. 'updatedBy' for log will be null.", id);
        }
        
        // Preserve original creation details if not meant to be updated
        // vehicle.setCreatedAt(existingVehicle.getCreatedAt());
        // vehicle.setCreatedBy(existingVehicle.getCreatedBy());
        // If Vehicle entity has 'updatedAt' field:
        // vehicle.setUpdatedAt(new Date());

        Vehicle updatedVehicle = service.save(vehicle);

        if (updatedBy != null) {
            Log_vehicle logEntry = new Log_vehicle(); // Ensure correct log model
            logEntry.setVehicleId(id); // Assuming this setter exists
            logEntry.setUpdatedAt(new Date());
            logEntry.setUpdatedBy(updatedBy);
            logVehicleRepository.save(logEntry);
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Vehicle updated successfully",
                "data", updatedVehicle
        ));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, Object>> delete(
            @RequestHeader(name = "X-User-ID", required = false) String userIdString,
            @RequestHeader(name = "X-User-Role", required = false) String userRole,
            @PathVariable Integer id) {

        log.info("Received delete vehicle request for id: {}. UserID: {}, Role: {}", id, userIdString, userRole);
        
        // Primary authorization is at the Gateway.
        service.deleteById(id); // Ensure service layer handles "not found" or check here
        
        if (userIdString != null && !userIdString.isEmpty()) {
            try {
                Integer deletedBy = Integer.parseInt(userIdString);
                Log_vehicle logEntry = new Log_vehicle();
                logEntry.setVehicleId(id); // Assuming setter exists
                logEntry.setUpdatedAt(new Date()); // Or a 'deletedAt' field
                logEntry.setUpdatedBy(deletedBy);  // Or 'deletedBy'
                logVehicleRepository.save(logEntry);
                log.info("Vehicle id: {} deleted by UserID: {}", id, deletedBy);
            } catch (NumberFormatException e) {
                log.warn("Could not parse X-User-ID for delete log for vehicle id: {}: {}", id, e.getMessage());
            }
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Vehicle deleted successfully",
                "data", ""
        ));
    }
}