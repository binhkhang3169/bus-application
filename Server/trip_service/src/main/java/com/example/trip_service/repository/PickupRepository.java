package com.example.trip_service.repository;

import com.example.trip_service.model.Pickup;
import com.example.trip_service.model.Province;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface PickupRepository extends JpaRepository<Pickup, String> {
    List<Pickup> findByStatusIn(List<Integer> statuses);
    List<Pickup> findByRouteIdAndSelfId(Integer routeId, String selfId);
}
