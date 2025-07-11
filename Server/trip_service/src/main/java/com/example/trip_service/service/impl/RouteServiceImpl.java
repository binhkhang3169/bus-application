package com.example.trip_service.service.impl;

import com.example.trip_service.model.Route;
import com.example.trip_service.model.SpecialDay;
import com.example.trip_service.model.log.Log_route;
import com.example.trip_service.repository.RouteRepository;
import com.example.trip_service.repository.log.LogRouteRepository;
import com.example.trip_service.service.intef.RouteService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class RouteServiceImpl implements RouteService {

    private final RouteRepository repository;

    @Autowired
    private LogRouteRepository logRouteRepository;

    public RouteServiceImpl(RouteRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Route> findAll() {
        return repository.findAll();
    }

    @Override
    public Route findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Route save(Route route) {
        if (route.getCreatedAt() == null) {
            route.setCreatedAt(new Date());
        }
        return repository.save(route);
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }

    @Override
    public List<Route> getRoutesByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updateRouteStatus(Integer userId,Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Route route = repository.findById(id).orElse(null);
        if (route!=null) {
            Route routeUpdate = route;
            routeUpdate.setStatus(newStatus);
            repository.save(routeUpdate);

            Log_route logRoute = new Log_route();
            logRoute.setRouteId(routeUpdate.getId());
            logRoute.setUpdatedAt(new Date());
            logRoute.setUpdatedBy(userId);
            logRouteRepository.save(logRoute);
            return true;
        }
        return false;
    }
}
