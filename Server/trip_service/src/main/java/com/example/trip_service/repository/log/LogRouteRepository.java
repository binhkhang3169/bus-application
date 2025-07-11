package com.example.trip_service.repository.log;

import com.example.trip_service.model.log.Log_route;
import com.example.trip_service.model.log.Log_specialDay;
import org.springframework.data.jpa.repository.JpaRepository;


public interface LogRouteRepository extends JpaRepository<Log_route, Integer> {
}

