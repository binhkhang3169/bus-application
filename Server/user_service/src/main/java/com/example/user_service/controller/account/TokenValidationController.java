package com.example.user_service.controller.account;

import com.example.user_service.dto.LoginRequest;
import com.example.user_service.model.RefreshToken;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.RefreshTokenService;
import com.example.user_service.service.UserService;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.JwtException;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Arrays;
import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1/auth")
public class TokenValidationController {

    @Autowired
    private JwtTokenUtil jwtTokenUtil;

    @PostMapping("/check")
    public ResponseEntity<Map<String, Object>> checkToken(@RequestHeader("Authorization") String token) {
        try {
            String role = jwtTokenUtil.getRoleFromToken(token);

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Token hợp lệ",
                    "data", Map.of("role", role)
            ));
        } catch (ExpiredJwtException e) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token đã hết hạn",
                    "data", ""
            ));
        } catch (JwtException e) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token giả mạo hoặc không hợp lệ",
                    "data", ""
            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi không xác định: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

}