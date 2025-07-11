package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_pickup;
import com.example.trip_service.model.log.Log_province;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogPickupRepository extends JpaRepository<Log_pickup, Integer> {
}

