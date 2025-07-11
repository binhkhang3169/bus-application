package com.example.trip_service.service.impl;

import com.example.trip_service.model.Province;
import com.example.trip_service.model.Route;
import com.example.trip_service.model.log.Log_province;
import com.example.trip_service.repository.ProvinceRepository;
import com.example.trip_service.repository.log.LogProvinceRepository;
import com.example.trip_service.service.intef.ProvinceService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class ProvinceServiceImpl implements ProvinceService {

    private final ProvinceRepository repository;

    @Autowired
    private LogProvinceRepository logProvinceRepository;

    public ProvinceServiceImpl(ProvinceRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Province> findAll() {
        return repository.findAll();
    }

    @Override
    public Province findById(Integer id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Province save(Province province) {
        if (province.getCreatedAt() == null) {
            province.setCreatedAt(new Date());
        }
        return repository.save(province);
    }

    @Override
    public void deleteById(Integer id) {
        repository.deleteById(id);
    }



    @Override
    public List<Province> getProvincesByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updateProvinceStatus(Integer userId,Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Province province = repository.findById(id).orElse(null);
        if (province!=null) {
            Province provinceUpdate = province;
            provinceUpdate.setStatus(newStatus);
            repository.save(provinceUpdate);

            Log_province logProvince = new Log_province();
            logProvince.setProvinceId(provinceUpdate.getId());
            logProvince.setUpdatedAt(new Date());
            logProvince.setUpdatedBy(userId);
            logProvinceRepository.save(logProvince);
            return true;
        }
        return false;
    }
}
