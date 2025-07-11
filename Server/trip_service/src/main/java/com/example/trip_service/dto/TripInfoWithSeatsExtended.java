package com.example.trip_service.dto;

import java.util.ArrayList;
import java.util.List;

public class TripInfoWithSeatsExtended extends TripInfoWithSeats {
    private List<Integer> availableSeatIds;
    
    public TripInfoWithSeatsExtended(TripInfoProjection tripInfo, List<String> availableSeats) {
        super(tripInfo, availableSeats);
        this.availableSeatIds = new ArrayList<>();
    }
    
    public List<Integer> getAvailableSeatIds() {
        return availableSeatIds;
    }
    
    public void setAvailableSeatIds(List<Integer> availableSeatIds) {
        this.availableSeatIds = availableSeatIds;
    }
}