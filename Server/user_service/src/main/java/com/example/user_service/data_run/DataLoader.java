package com.example.user_service.data_run;

import com.example.user_service.model.Role;
import com.example.user_service.model.User;
import com.example.user_service.repository.RoleRepository;
import com.example.user_service.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Collections;
import java.util.Date;

@Configuration
public class DataLoader {

    @Autowired
    private RoleRepository roleRepository;

    private PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();


    @Autowired
    private UserRepository userRepository;

    @Bean
    public CommandLineRunner loadData() {
        return args -> {
            if (roleRepository.findByRoleName("ROLE_ADMIN") == null) {
                roleRepository.save(new Role("ROLE_ADMIN"));
            }

            if (roleRepository.findByRoleName("ROLE_CUSTOMER") == null) {
                roleRepository.save(new Role("ROLE_CUSTOMER"));
            }

            if (roleRepository.findByRoleName("ROLE_DRIVER") == null) {
                roleRepository.save(new Role("ROLE_DRIVER"));
            }


            if (roleRepository.findByRoleName("ROLE_RECEPTION") == null) {
                roleRepository.save(new Role("ROLE_RECEPTION"));
            }

            if (roleRepository.findByRoleName("ROLE_OPERATOR") == null) {
                roleRepository.save(new Role("ROLE_OPERATOR"));
            }

            String adminEmail = "caoky.sonha@gmail.com";

            if(userRepository.findByUsername(adminEmail)==null){
                User adminUser = new User();
                adminUser.setCreatedAt(new Date());
                adminUser.setUsername(adminEmail);
                adminUser.setPassword(passwordEncoder.encode("nhaxeanhphung"));
                Role adminRole = roleRepository.findByRoleName("ROLE_ADMIN");
                adminUser.setRoles(Collections.singleton(adminRole));
                userRepository.save(adminUser);
                System.out.println("Created admin user: " + adminEmail + " with password: nhaxeanhphung");

            }


        };
    }


}