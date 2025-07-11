package com.example.trip_service.service.intef;

import com.example.trip_service.model.Pickup;
import com.example.trip_service.model.Province;

import java.util.List;

public interface PickupService {
    List<Pickup> findAll();
    Pickup findById(String id);
    Pickup save(Pickup pickup);
    void deleteById(String id);
    public List<Pickup> getPickupsByStatus(Integer status) ;

    public boolean updatePickupStatus(Integer userId,String id, Integer newStatus) ;
    List<Pickup> findByRouteIdAndSelfIdIsMinusOne(Integer routeId);

}
