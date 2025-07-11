package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_province;
import com.example.trip_service.model.log.Log_route;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogProvinceRepository extends JpaRepository<Log_province, Integer> {
}

