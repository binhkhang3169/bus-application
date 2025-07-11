package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_trip;
import com.example.trip_service.model.log.Log_type;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogTripRepository extends JpaRepository<Log_trip, Integer> {
}

