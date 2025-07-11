package com.example.trip_service.service.impl;

import com.example.trip_service.model.SpecialDay;
import com.example.trip_service.model.Station;
import com.example.trip_service.model.log.Log_specialDay;
import com.example.trip_service.repository.SpecialDayRepository;
import com.example.trip_service.repository.log.LogSpecialDayRepository;
import com.example.trip_service.service.intef.SpecialDayService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class SpecialDayServiceImpl implements SpecialDayService {

    private final SpecialDayRepository repository;

    @Autowired
    private LogSpecialDayRepository logSpecialDayRepository;

    public SpecialDayServiceImpl(SpecialDayRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<SpecialDay> findAll() {
        return repository.findAll();
    }

    @Override
    public SpecialDay findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public SpecialDay save(SpecialDay day) {
        if (day.getCreatedAt() == null) {
            day.setCreatedAt(new Date());
        }
        return repository.save(day);
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }



    @Override
    public List<SpecialDay> getSpecialDaysByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updateSpecialDayStatus(Integer userId,Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        SpecialDay specialDay = repository.findById(id).orElse(null);
        if (specialDay!=null) {
            SpecialDay specialDayUpdate = specialDay;
            specialDayUpdate.setStatus(newStatus);
            repository.save(specialDayUpdate);

            Log_specialDay logSpecialDay = new Log_specialDay();
            logSpecialDay.setSpecialDayId(specialDayUpdate.getId());
            logSpecialDay.setUpdatedAt(new Date());
            logSpecialDay.setUpdatedBy(userId);
            logSpecialDayRepository.save(logSpecialDay);
            return true;
        }
        return false;
    }

}
