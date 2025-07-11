package com.example.trip_service.repository;

import com.example.trip_service.model.Station;
import com.example.trip_service.model.Trip;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface StationRepository extends JpaRepository<Station, Integer> {
    List<Station> findByStatusIn(List<Integer> statuses);

}
