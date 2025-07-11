package com.example.trip_service.dto;


public interface TripInfoProjection {
    Integer getTripId();
    String getVehicleId();
    String getLicense();
    String getVehicleType();


    String getStatus();
    Integer getStock();
    String getDepartureDate();
    String getDepartureTime();
    String getArrivalDate();
    String getArrivalTime();
    Integer getPrice();
    String getEstimatedDistance();
    String getEstimatedTime();
    String getDepartureStation();
    String getArrivalStation();
    String getFullRoute();
}
