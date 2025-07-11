package com.example.trip_service.dto;

import com.example.trip_service.model.Seat;

import java.util.ArrayList;
import java.util.List;

public class TripInfoWithDetailedSeats extends TripInfoWithSeats {
    private List<Seat> detailedSeats;
    
    public TripInfoWithDetailedSeats(TripInfoProjection tripInfo, List<String> availableSeats) {
        super(tripInfo, availableSeats);
        this.detailedSeats = new ArrayList<>();
    }
    
    public List<Seat> getDetailedSeats() {
        return detailedSeats;
    }
    
    public void setDetailedSeats(List<Seat> detailedSeats) {
        this.detailedSeats = detailedSeats;
    }
}