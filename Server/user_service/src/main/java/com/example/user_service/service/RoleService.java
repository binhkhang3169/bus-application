package com.example.user_service.service;


import com.example.user_service.model.Role;
import com.example.user_service.repository.RoleRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class RoleService {
    @Autowired
    private RoleRepository roleRepository;

    public Role createRole(String roleName) {
        Role role = new Role(roleName);
        return roleRepository.save(role);
    }
    public RoleService(RoleRepository roleRepository) {
        this.roleRepository = roleRepository;
    }
    public Role findByName(String roleName) {
        return roleRepository.findByRoleName(roleName);
    }
}
