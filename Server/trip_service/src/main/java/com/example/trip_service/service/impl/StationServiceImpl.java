package com.example.trip_service.service.impl;

import com.example.trip_service.model.Station;
import com.example.trip_service.model.Trip;
import com.example.trip_service.model.log.Log_station;
import com.example.trip_service.repository.StationRepository;
import com.example.trip_service.repository.log.LogStationRepository;
import com.example.trip_service.service.intef.StationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class StationServiceImpl implements StationService {

    private final StationRepository repository;

    @Autowired
    private LogStationRepository logStationRepository;

    public StationServiceImpl(StationRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Station> findAll() {
        return repository.findAll();
    }

    @Override
    public Station findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Station save(Station station) {
        if (station.getCreatedAt() == null) {
            station.setCreatedAt(new Date());
        }
        return repository.save(station);
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }

    @Override
    public List<Station> getStationsByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }

    @Override
    public boolean updateStationStatus(Integer userId, Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Station station = repository.findById(id).orElse(null);
        if (station != null) {
            Station stationUpdate = station;
            stationUpdate.setStatus(newStatus);
            repository.save(stationUpdate);

            Log_station logStation = new Log_station();
            logStation.setStationId(stationUpdate.getId());
            logStation.setUpdatedAt(new Date());
            logStation.setUpdatedBy(userId);
            logStationRepository.save(logStation);
            return true;
        }
        return false;
    }


}
