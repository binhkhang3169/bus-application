package com.example.user_service.controller.account;

import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.EmailService;
import com.example.user_service.service.UserService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.ModelAndView;

import java.util.Map;
import java.util.UUID;

import static com.example.user_service.domain.IP.IP_network;

@RestController
@RequestMapping("/api/v1")
public class PasswordResetController {

    private final UserService userService;
    private final EmailService emailService;
    private final PasswordEncoder passwordEncoder = new BCryptPasswordEncoder(); // Made final

    public PasswordResetController(UserService userService, EmailService emailService) {
        this.userService = userService;
        this.emailService = emailService;
    }

    @PostMapping("/change-password")
    public ResponseEntity<Map<String, Object>> changePassword( // ResponseEntity updated
            @RequestHeader("Authorization") String token,
            @RequestParam String oldPassword,
            @RequestParam String newPassword) {
        try {
            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
            User user = userService.getUserById(userId);

            if (user == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy người dùng",
                        "data", ""
                ));
            }

            if (!passwordEncoder.matches(oldPassword, user.getPassword())) {
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                        "code", HttpStatus.BAD_REQUEST.value(),
                        "message", "Mật khẩu cũ không đúng",
                        "data", ""
                ));
            }

            user.setPassword(passwordEncoder.encode(newPassword));
            userService.save(user);

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Đổi mật khẩu thành công",
                    "data", ""
            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi đổi mật khẩu: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    @PostMapping("/forgot-password")
    public ResponseEntity<Map<String, Object>> forgotPassword(@RequestParam String email) { // ResponseEntity updated
        User user = userService.findByUsername(email);
        if (user == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                    "code", HttpStatus.NOT_FOUND.value(),
                    "message", "Email không tồn tại trong hệ thống",
                    "data", ""
            ));
        }

        try {
            String token = UUID.randomUUID().toString();
            userService.savePasswordResetToken(user, token);

            String resetLink = "http://" + IP_network + ":8080/api/v1/reset-password?token=" + token; // Corrected path to match GET mapping

            emailService.sendEmail(email, "Đặt lại mật khẩu", "Nhấp vào liên kết sau để đặt lại mật khẩu: " + resetLink);

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Email đặt lại mật khẩu đã được gửi",
                    "data", ""
            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi gửi email đặt lại mật khẩu: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    /**
     * This endpoint serves an HTML page for password reset and is not a typical JSON API.
     * Therefore, it will continue to return ModelAndView and is not refactored to the standard JSON response.
     */
    @GetMapping("/reset-password")
    public ModelAndView showResetPasswordPage(@RequestParam String token) {
        ModelAndView modelAndView = new ModelAndView("reset_password"); // Standard view name
        boolean isValid = userService.validateResetToken(token);

        if (!isValid) {
            modelAndView.setViewName("error"); // Error view
            modelAndView.addObject("message", "Token không hợp lệ hoặc đã hết hạn.");
            return modelAndView;
        }

        modelAndView.addObject("token", token);
        return modelAndView;
    }

    @PostMapping("/reset-password")
    public ResponseEntity<Map<String, Object>> processResetPassword(@RequestParam String token, @RequestParam String newPassword) { // ResponseEntity updated
        boolean isUpdated = userService.updatePassword(token, newPassword);

        if (!isUpdated) {
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Token không hợp lệ hoặc đã hết hạn.",
                    "data", ""
            ));
        }

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Mật khẩu đã được cập nhật thành công.",
                "data", ""
        ));
    }
}