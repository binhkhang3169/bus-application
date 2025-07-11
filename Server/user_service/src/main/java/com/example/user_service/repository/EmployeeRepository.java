package com.example.user_service.repository;

import com.example.user_service.model.Customer;
import com.example.user_service.model.Employee;
import org.springframework.data.jpa.repository.JpaRepository;

public interface EmployeeRepository extends JpaRepository<Employee, Integer> {



}

