package com.example.trip_service.service.impl;

import com.example.trip_service.model.Pickup;
import com.example.trip_service.model.Province;
import com.example.trip_service.model.log.Log_pickup;
import com.example.trip_service.repository.PickupRepository;
import com.example.trip_service.repository.log.LogPickupRepository;
import com.example.trip_service.service.intef.PickupService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class PickupServiceImpl implements PickupService {

    private final PickupRepository repository;

    @Autowired
    private LogPickupRepository logPickupRepository;

    public PickupServiceImpl(PickupRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Pickup> findAll() {
        return repository.findAll();
    }

    @Override
    public Pickup findById(String id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Pickup save(Pickup pickup) {
        if (pickup.getCreatedAt() == null) {
            pickup.setCreatedAt(new Date());
        }
        return repository.save(pickup);
    }

    @Override
    public void deleteById(String id) {
        repository.deleteById(id);
    }



    @Override
    public List<Pickup> getPickupsByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updatePickupStatus(Integer userId,String id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Pickup pickup = repository.findById(id).orElse(null);
        if (pickup!=null) {
            Pickup pickupUpdate = pickup;
            pickupUpdate.setStatus(newStatus);
            repository.save(pickupUpdate);

            Log_pickup logPickup = new Log_pickup();
            logPickup.setPickupId(pickupUpdate.getId());
            logPickup.setUpdatedAt(new Date());
            logPickup.setUpdatedBy(userId);
            logPickupRepository.save(logPickup);
            return true;
        }
        return false;
    }

    @Override
    public List<Pickup> findByRouteIdAndSelfIdIsMinusOne(Integer routeId) {
        if (routeId == null) {
            // Or throw an IllegalArgumentException, depending on how you want to handle invalid input
            return List.of(); // Return an empty list if routeId is null
        }
        return repository.findByRouteIdAndSelfId(routeId, "-1");
    }
}


