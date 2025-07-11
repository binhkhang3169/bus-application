package com.example.user_service.repository;

import com.example.user_service.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.Optional;

public interface UserRepository extends JpaRepository<User, Integer> {


    Optional<User> findUserByUsername(String username);
    User findByUsername(String username);


//    @Query("SELECT u.id FROM Users u WHERE u.email = :email")
//    Optional<Integer> findIdByEmail(@Param("email") String email);

    boolean existsByUsername(String username);

    User findByResetToken(String resetToken);
    
}

