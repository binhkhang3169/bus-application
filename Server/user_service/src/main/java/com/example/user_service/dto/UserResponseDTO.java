package com.example.user_service.dto;

import com.example.user_service.model.Customer;
import com.example.user_service.model.Employee;
import com.example.user_service.model.Role;
import com.example.user_service.model.User;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class UserResponseDTO {
    private int id;
    private String username;
    private Date createdAt;
    private Integer active;
    private Set<String> roles;
    private Map<String, Object> details;

    // Empty constructor
    public UserResponseDTO() {
    }

    // Constructor from User entity
    public UserResponseDTO(User user) {
        this.id = user.getId();
        this.username = user.getUsername();
        this.createdAt = user.getCreatedAt();
        this.active = user.getActive();
        this.roles = user.getRoles().stream()
                .map(Role::getRoleName)
                .collect(Collectors.toSet());
        this.details = new HashMap<>();
    }

    // Static factory method to create from User with Customer details
    public static UserResponseDTO fromCustomer(User user, Customer customer) {
        UserResponseDTO dto = new UserResponseDTO(user);
        if (customer != null) {
            Map<String, Object> customerDetails = new HashMap<>();
            customerDetails.put("fullName", customer.getFullName());
            customerDetails.put("phoneNumber", customer.getPhoneNumber());
            customerDetails.put("address", customer.getAddress());
            customerDetails.put("gender", customer.getGender());
            dto.setDetails(customerDetails);
        }
        return dto;
    }

    // Static factory method to create from User with Employee details
    public static UserResponseDTO fromEmployee(User user, Employee employee) {
        UserResponseDTO dto = new UserResponseDTO(user);
        if (employee != null) {
            Map<String, Object> employeeDetails = new HashMap<>();
            employeeDetails.put("fullName", employee.getFullName());
            employeeDetails.put("phoneNumber", employee.getPhoneNumber());
            employeeDetails.put("address", employee.getAddress());
            employeeDetails.put("gender", employee.getGender());
            employeeDetails.put("dateOfBirth", employee.getDateOfBirth());
            employeeDetails.put("hiredDate", employee.getHiredDate());
            employeeDetails.put("identityNumber", employee.getIdentityNumber());
            employeeDetails.put("avatarUrl", employee.getAvatarUrl());
            dto.setDetails(employeeDetails);
        }
        return dto;
    }

    // Getters and Setters
    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }

    public Integer getActive() {
        return active;
    }

    public void setActive(Integer active) {
        this.active = active;
    }

    public Set<String> getRoles() {
        return roles;
    }

    public void setRoles(Set<String> roles) {
        this.roles = roles;
    }

    public Map<String, Object> getDetails() {
        return details;
    }

    public void setDetails(Map<String, Object> details) {
        this.details = details;
    }
}