package com.example.trip_service.service.intef;

import com.example.trip_service.model.Province;
import com.example.trip_service.model.Route;

import java.util.List;

public interface ProvinceService {
    List<Province> findAll();
    Province findById(Integer id);
    Province save(Province province);
    void deleteById(Integer id);

    public List<Province> getProvincesByStatus(Integer status) ;

    public boolean updateProvinceStatus(Integer userId,Integer id, Integer newStatus) ;
}
