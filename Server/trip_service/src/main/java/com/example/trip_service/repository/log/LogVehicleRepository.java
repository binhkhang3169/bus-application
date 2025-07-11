package com.example.trip_service.repository.log;

import com.example.trip_service.model.Vehicle;
import com.example.trip_service.model.log.Log_vehicle;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;


public interface LogVehicleRepository extends JpaRepository<Log_vehicle, Integer> {
}

