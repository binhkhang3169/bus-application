package com.example.trip_service.dto;

import java.util.Date;

public record TripCreatedEvent(
    String tripId,
    Integer totalSeats,
    Date creationTimestamp
) {}