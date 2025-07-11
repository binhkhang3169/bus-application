package com.example.trip_service.dto;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class SeatUpdateEvent {
    private String tripId;
    private int seatCount;
}