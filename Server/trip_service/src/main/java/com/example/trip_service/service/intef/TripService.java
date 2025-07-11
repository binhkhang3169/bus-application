package com.example.trip_service.service.intef;

import java.util.List;

import com.example.trip_service.dto.TripInfoProjection;
import com.example.trip_service.model.Trip;

public interface TripService {
    List<Trip> findAll();
    Trip findById(int id);
    Trip save(Trip trip);
    void deleteById(int id);

    // Original methods that return TripInfoProjection
    List<TripInfoProjection> getTripInfos(Integer routeId, String departureDate, Integer quantity);
    List<TripInfoProjection> searchTripsByLocations(Integer fromProvinceId, Integer toProvinceId, String departureDate, Integer quantity);

    TripInfoProjection findByIdWithSeats(int id);

    // New methods that return TripInfoWithSeats
    List<TripInfoProjection> searchTripsByLocationsWithSeats(Integer fromProvinceId, Integer toProvinceId, String departureDate, Integer quantity, Integer userId);



    List<Trip> getTripsByStatus(Integer status);
    boolean updateTripStatus(Integer userId, Integer id, Integer newStatus);

    // New service method
    List<Trip> findTripsByDriverId(Integer driverId);
}