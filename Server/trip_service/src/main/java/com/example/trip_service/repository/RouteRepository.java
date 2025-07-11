package com.example.trip_service.repository;

import com.example.trip_service.model.Route;
import com.example.trip_service.model.SpecialDay;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface RouteRepository extends JpaRepository<Route, Integer> {
    List<Route> findByStatusIn(List<Integer> statuses);

}
