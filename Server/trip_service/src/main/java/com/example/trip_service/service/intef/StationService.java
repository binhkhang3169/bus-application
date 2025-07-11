package com.example.trip_service.service.intef;

import com.example.trip_service.model.Station;
import com.example.trip_service.model.Trip;

import java.util.List;

public interface StationService {
    List<Station> findAll();
    Station findById(int id);
    Station save(Station station);
    void deleteById(int id);

    public List<Station> getStationsByStatus(Integer status) ;

    public boolean updateStationStatus(Integer userId,Integer id, Integer newStatus) ;
}
