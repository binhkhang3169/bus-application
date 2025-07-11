package com.example.trip_service.service.intef;

import com.example.trip_service.model.SpecialDay;
import com.example.trip_service.model.Station;

import java.util.List;

public interface SpecialDayService {
    List<SpecialDay> findAll();
    SpecialDay findById(int id);
    SpecialDay save(SpecialDay day);
    void deleteById(int id);

    public List<SpecialDay> getSpecialDaysByStatus(Integer status) ;

    public boolean updateSpecialDayStatus(Integer userId,Integer id, Integer newStatus) ;
}
