package com.example.trip_service.repository;

import com.example.trip_service.model.Province;
import com.example.trip_service.model.Route;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ProvinceRepository extends JpaRepository<Province, Integer> {
    List<Province> findByStatusIn(List<Integer> statuses);

}
