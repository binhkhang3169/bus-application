package com.example.trip_service.dto;

import java.util.List;

public class TripInfoWithSeats {
    private TripInfoProjection tripInfo;
    private List<String> availableSeats;
    
    public TripInfoWithSeats(TripInfoProjection tripInfo, List<String> availableSeats) {
        this.tripInfo = tripInfo;
        this.availableSeats = availableSeats;
    }
    
    public TripInfoProjection getTripInfo() {
        return tripInfo;
    }
    
    public void setTripInfo(TripInfoProjection tripInfo) {
        this.tripInfo = tripInfo;
    }
    
    public List<String> getAvailableSeats() {
        return availableSeats;
    }
    
    public void setAvailableSeats(List<String> availableSeats) {
        this.availableSeats = availableSeats;
    }
}