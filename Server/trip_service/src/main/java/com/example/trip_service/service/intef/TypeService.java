package com.example.trip_service.service.intef;

import com.example.trip_service.model.Type;
import com.example.trip_service.model.Vehicle;

import java.util.List;

public interface TypeService {
    List<Type> findAll();
    Type findById(int id);
    Type save(Type vehicle);
    void deleteById(int id);

    public List<Type> getTypesByStatus(Integer status) ;

    public boolean updateTypeStatus(Integer userId,Integer id, Integer newStatus) ;

}
