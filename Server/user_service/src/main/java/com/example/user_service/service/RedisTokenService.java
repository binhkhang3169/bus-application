package com.example.user_service.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.Duration;

@Service
public class RedisTokenService {

    @Autowired
    private RedisTemplate<String, Object> redisTemplate;

    private final long TOKEN_TTL = 60 * 60; // 1 giờ

    public void saveRefreshToken(Integer userId, String refreshToken) {
        String key = "refresh_token:" + userId;
        redisTemplate.opsForValue().set(key, refreshToken, Duration.ofSeconds(TOKEN_TTL));
    }

    public String getRefreshToken(Integer userId) {
        String key = "refresh_token:" + userId;
        return (String) redisTemplate.opsForValue().get(key);
    }

    public void deleteRefreshToken(Integer userId) {
        redisTemplate.delete("refresh_token:" + userId);
    }
}
