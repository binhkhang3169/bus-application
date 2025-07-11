package com.example.trip_service.service.impl;

import com.example.trip_service.model.Vehicle;
import com.example.trip_service.model.log.Log_vehicle;
import com.example.trip_service.repository.VehicleRepository;
import com.example.trip_service.repository.log.LogVehicleRepository;
import com.example.trip_service.service.intef.VehicleService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class VehicleServiceImpl implements VehicleService {

    private static final Logger log = LoggerFactory.getLogger(VehicleServiceImpl.class);
    private final VehicleRepository repository;

    @Autowired
    private  LogVehicleRepository logVehicleRepository;



    public VehicleServiceImpl(VehicleRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Vehicle> findAll() {
        return repository.findAll();
    }

    @Override
    public Vehicle findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Vehicle save(Vehicle vehicle) {
        if (vehicle.getCreatedAt() == null) {
            vehicle.setCreatedAt(new Date());
        }
        return repository.save(vehicle);
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }

    @Override
    public List<Vehicle> getVehiclesByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updateVehicleStatus(Integer userId,Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Vehicle optionalVehicle = repository.findById(id).orElse(null);
        if (optionalVehicle!=null) {
            Vehicle vehicle = optionalVehicle;
            vehicle.setStatus(newStatus);
            repository.save(vehicle);
            Log_vehicle logVehicle = new Log_vehicle();
            logVehicle.setVehicleId(vehicle.getId());
            logVehicle.setUpdatedAt(new Date());
            logVehicle.setUpdatedBy(userId);
            logVehicleRepository.save(logVehicle);
            return true;
        }
        return false;
    }





}
