package com.example.trip_service.repository;

import com.example.trip_service.model.Vehicle;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface VehicleRepository extends JpaRepository<Vehicle, Integer> {
    List<Vehicle> findByStatusIn(List<Integer> statuses);
}
