package com.example.user_service.service;


import com.example.user_service.model.Customer;
import com.example.user_service.model.Employee;
import com.example.user_service.repository.CustomerRepository;
import com.example.user_service.repository.EmployeeRepository;
import com.example.user_service.repository.RoleRepository;
import com.example.user_service.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class EmployeeService {
    @Autowired
    RoleRepository roleRepository;


    @Autowired
    private UserRepository userRepository;

    @Autowired
    private EmployeeRepository employeeRepository;

    public void save(Employee employee) {
        employeeRepository.save(employee);
    }


    public void updateUser(Employee employee) {
        employeeRepository.save(employee);
    }


    public Employee findById(int id) {
        return employeeRepository.findById(id).get();
    }


}