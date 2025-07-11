package com.example.trip_service.controller;

import com.example.trip_service.util.LocationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1/locations")
public class LocationController {

    private final LocationService locationService;

    @Autowired
    public LocationController(LocationService locationService) {
        this.locationService = locationService;
    }

    @GetMapping
    public ResponseEntity<Map<String, Object>> getAllLocations() {
        List<Map<String, Object>> locations = locationService.getAllLocations();
        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Locations retrieved successfully",
                "data", locations
        ));
    }

    @GetMapping("/code/{code}")
    public ResponseEntity<Map<String, Object>> getLocationByCode(@PathVariable String code) {
        Optional<Map<String, Object>> locationOpt = locationService.findLocationByCode(code);
        if (locationOpt.isPresent()) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Location found by code",
                    "data", locationOpt.get()
            ));
        } else {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Location not found by code",
                    "data", ""
            ));
        }
    }

    @GetMapping("/id/{id}")
    public ResponseEntity<Map<String, Object>> getLocationById(@PathVariable Integer id) {
        Optional<Map<String, Object>> locationOpt = locationService.findLocationById(id);
        if (locationOpt.isPresent()) {
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Location found by ID",
                    "data", locationOpt.get()
            ));
        } else {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Location not found by ID",
                    "data", ""
            ));
        }
    }
}