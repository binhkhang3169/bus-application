package com.example.user_service.service;


import com.example.user_service.model.Customer;
import com.example.user_service.repository.CustomerRepository;
import com.example.user_service.repository.RoleRepository;
import com.example.user_service.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class CustomerService {
    @Autowired
    RoleRepository roleRepository;


    @Autowired
    private UserRepository userRepository;

    @Autowired
    private CustomerRepository customerRepository;

    public void save(Customer customer) {
        customerRepository.save(customer);
    }


    public void updateUser(Customer customer) {
        customerRepository.save(customer);
    }


    public Customer findById(int id) {
        return customerRepository.findById(id).get();
    }


}