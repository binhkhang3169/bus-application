package com.example.trip_service.service.intef;

import com.example.trip_service.model.Route;
import com.example.trip_service.model.SpecialDay;

import java.util.List;

public interface RouteService {
    List<Route> findAll();
    Route findById(int id);
    Route save(Route route);
    void deleteById(int id);

    public List<Route> getRoutesByStatus(Integer status) ;

    public boolean updateRouteStatus(Integer userId,Integer id, Integer newStatus) ;
}
