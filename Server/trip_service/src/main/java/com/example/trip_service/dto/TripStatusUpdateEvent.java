package com.example.trip_service.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class TripStatusUpdateEvent {
    private Integer tripId;
    private Integer oldStatus;
    private Integer newStatus;
    private Integer updatedBy;
    private Date updatedAt;
}