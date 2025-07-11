package com.example.trip_service.util;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.core.type.TypeReference;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

@Component
public class LocationService {
    private List<Map<String, Object>> locations;
    private Map<String, Integer> codeToIdMap;
    private Map<Integer, String> idToCodeMap;

    @PostConstruct
    public void init() throws IOException {
        try {
            ObjectMapper mapper = new ObjectMapper();
            InputStream is = new ClassPathResource("locations.json").getInputStream();
            locations = mapper.readValue(is, new TypeReference<List<Map<String, Object>>>() {});
            
            // Initialize maps for quick lookups
            codeToIdMap = locations.stream()
                .collect(Collectors.toMap(
                    location -> (String) location.get("code"),
                    location -> ((Number) location.get("id")).intValue()
                ));
                
            idToCodeMap = locations.stream()
                .collect(Collectors.toMap(
                    location -> ((Number) location.get("id")).intValue(),
                    location -> (String) location.get("code")
                ));
        } catch (IOException e) {
            throw new RuntimeException("Failed to load locations data", e);
        }
    }
    
    /**
     * Get location ID by code
     */
    public Integer getLocationIdByCode(String code) {
        return codeToIdMap.get(code);
    }
    
    /**
     * Get location code by ID
     */
    public String getLocationCodeById(Integer id) {
        return idToCodeMap.get(id);
    }
    
    /**
     * Get all locations
     */
    public List<Map<String, Object>> getAllLocations() {
        return locations;
    }
    
    /**
     * Find location by ID
     */
    public Optional<Map<String, Object>> findLocationById(Integer id) {
        return locations.stream()
            .filter(location -> ((Number) location.get("id")).intValue() == id)
            .findFirst();
    }
    
    /**
     * Find location by code
     */
    public Optional<Map<String, Object>> findLocationByCode(String code) {
        return locations.stream()
            .filter(location -> location.get("code").equals(code))
            .findFirst();
    }
}