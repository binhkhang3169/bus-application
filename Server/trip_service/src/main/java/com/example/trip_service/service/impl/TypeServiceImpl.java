package com.example.trip_service.service.impl;

import com.example.trip_service.model.Type;
import com.example.trip_service.model.Vehicle;
import com.example.trip_service.model.log.Log_type;
import com.example.trip_service.repository.TypeRepository;
import com.example.trip_service.repository.log.LogTypeRepository;
import com.example.trip_service.service.intef.TypeService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

@Service
public class TypeServiceImpl implements TypeService {

    private static final Logger log = LoggerFactory.getLogger(TypeServiceImpl.class);
    private final TypeRepository repository;
    @Autowired
    private LogTypeRepository logTypeRepository;

    public TypeServiceImpl(TypeRepository repository) {
        this.repository = repository;
    }

    @Override
    public List<Type> findAll() {
        return repository.findAll();
    }

    @Override
    public Type findById(int id) {
        return repository.findById(id).orElse(null);
    }

    @Override
    public Type save(Type type) {
        if (type.getCreatedAt() == null) {
            type.setCreatedAt(new Date());
        }
        return repository.save(type);
    }

    @Override
    public void deleteById(int id) {
        repository.deleteById(id);
    }

    @Override
    public List<Type> getTypesByStatus(Integer status) {
        if (status != 0 && status != 1) {
            throw new IllegalArgumentException("Chỉ chấp nhận status là 0 hoặc 1");
        }
        return repository.findByStatusIn(List.of(status));
    }
    @Override
    public boolean updateTypeStatus(Integer userId,Integer id, Integer newStatus) {
        if (newStatus != 0 && newStatus != 1) {
            throw new IllegalArgumentException("Trạng thái không hợp lệ! Chỉ chấp nhận 0 hoặc 1.");
        }
        Type type = repository.findById(id).orElse(null);
        if (type!=null) {
            Type typeUpdate = type;
            typeUpdate.setStatus(newStatus);
            repository.save(typeUpdate);

            Log_type logType = new Log_type();
            logType.setTypeId(type.getId());
            logType.setUpdatedAt(new Date());
            logType.setUpdatedBy(userId);
            logTypeRepository.save(logType);
            return true;
        }
        return false;
    }
}
