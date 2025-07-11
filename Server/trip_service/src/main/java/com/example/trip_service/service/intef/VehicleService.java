package com.example.trip_service.service.intef;

import com.example.trip_service.model.Vehicle;

import java.util.List;

public interface VehicleService {
    List<Vehicle> findAll();
    Vehicle findById(int id);
    Vehicle save(Vehicle vehicle);
    void deleteById(int id);

    public List<Vehicle> getVehiclesByStatus(Integer status) ;

    public boolean updateVehicleStatus(Integer userId, Integer id, Integer newStatus) ;

}
