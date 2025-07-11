package com.example.trip_service.repository;

import com.example.trip_service.model.SpecialDay;
import com.example.trip_service.model.Station;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface SpecialDayRepository extends JpaRepository<SpecialDay, Integer> {
    List<SpecialDay> findByStatusIn(List<Integer> statuses);


}
