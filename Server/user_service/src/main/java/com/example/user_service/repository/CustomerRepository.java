package com.example.user_service.repository;

import com.example.user_service.model.Customer;
import com.example.user_service.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface CustomerRepository extends JpaRepository<Customer, Integer> {



}

