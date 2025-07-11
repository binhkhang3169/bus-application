package com.example.trip_service.repository;

import com.example.trip_service.model.Type;
import com.example.trip_service.model.Vehicle;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface TypeRepository extends JpaRepository<Type, Integer> {
    List<Type> findByStatusIn(List<Integer> statuses);

}
