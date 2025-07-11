package com.example.user_service.repository;

import com.example.user_service.model.RefreshToken;
import com.example.user_service.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface RefreshTokenRepository extends JpaRepository<RefreshToken, Long> {
    Optional<RefreshToken> findByToken(String token);

    @Modifying
    int deleteByUser(User user);

    @Modifying
    void deleteByToken(String token); // Added for direct deletion by token string
}