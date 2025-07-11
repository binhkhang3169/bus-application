package com.example.user_service.security;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

import java.util.Arrays;

@Configuration
public class SecurityConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }


    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
                // Disable CSRF
                .csrf().disable()
                .authorizeHttpRequests(authorize -> authorize
                        .requestMatchers("/**").permitAll()
                        // Your commented out role-based authorizations
                        // .requestMatchers("/admin","/manager-customer","/manager-package/**").hasRole("ADMIN")
                        // .requestMatchers("/customer","/customer/deposits","/customer-history","/customer-package","/customer-package/purchase","/customer-dangtin","/customer-history-list","/post").hasRole("CUSTOMER")
                        // .anyRequest().authenticated()
                )
                .addFilterBefore(new JwtAuthenticationFilter(), UsernamePasswordAuthenticationFilter.class)
                .logout(logout -> logout.permitAll());
        
        return http.build();
    }
}