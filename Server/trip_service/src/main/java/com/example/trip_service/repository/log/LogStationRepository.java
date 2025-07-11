package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_station;
import com.example.trip_service.model.log.Log_trip;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogStationRepository extends JpaRepository<Log_station, Integer> {
}

