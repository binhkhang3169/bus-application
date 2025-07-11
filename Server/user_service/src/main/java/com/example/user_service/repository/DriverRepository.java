package com.example.user_service.repository;

import com.example.user_service.model.DriverInfo;
import com.example.user_service.model.Employee;
import org.springframework.data.jpa.repository.JpaRepository;

public interface DriverRepository extends JpaRepository<DriverInfo, Integer> {



}

