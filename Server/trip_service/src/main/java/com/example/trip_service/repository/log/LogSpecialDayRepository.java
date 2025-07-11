package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_specialDay;
import com.example.trip_service.model.log.Log_station;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogSpecialDayRepository extends JpaRepository<Log_specialDay, Integer> {
}

