package com.example.user_service.controller.account;

import com.example.user_service.dto.SignupCustomerRequest;
import com.example.user_service.service.EmailService;
import com.example.user_service.service.RedisTokenService;
import com.example.user_service.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;
import java.util.Random;

@RestController
@RequestMapping("/api/v1")
public class SignupCustomerController {

    private final UserService userService; // Made final
    private final EmailService emailService; // Made final

    @Autowired
    public SignupCustomerController(UserService userService, EmailService emailService) {
        this.userService = userService;
        this.emailService = emailService;
    }

    @PostMapping("/signup")
    public ResponseEntity<Map<String, Object>> register(@RequestBody SignupCustomerRequest registerRequest) { // ResponseEntity updated
        if (userService.existsByUsername(registerRequest.getUsername())) {
            return ResponseEntity.badRequest().body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Email đã được sử dụng",
                    "data", ""
            ));
        }

        try {
            String otp = generateOTP();
            userService.saveOtp(registerRequest.getUsername(), otp);
            emailService.sendOtpEmail(registerRequest.getUsername(), otp);

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "OTP đã được gửi đến email. Vui lòng xác nhận!",
                    "data", ""
            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi trong quá trình đăng ký: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    @PostMapping("/verify-otp")
    public ResponseEntity<Map<String, Object>> verifyOtp(@RequestBody SignupCustomerRequest validRequest) { // ResponseEntity updated
        boolean isValid = userService.verifyOtp(validRequest.getUsername(), validRequest.getOtp());
        if (!isValid) {
            return ResponseEntity.badRequest().body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "OTP không hợp lệ hoặc đã hết hạn",
                    "data", ""
            ));
        }
        try {
            userService.registerUser(validRequest);
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Xác thực tài khoản thành công",
                    "data", ""
            ));
        } catch (Exception e) {
             return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi đăng ký người dùng sau xác thực OTP: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    @PostMapping("/resend-otp")
    public ResponseEntity<Map<String, Object>> resendOtp(@RequestBody Map<String, String> request) { // ResponseEntity updated
        String email = request.get("email");
        // Note: Original logic checks if email exists and returns bad request if it does.
        // This might be counter-intuitive for "resend OTP" if the user *has* an account but didn't get OTP.
        // Assuming the original logic is intended:
        if (userService.existsByUsername(email)) {
             // If OTP is for an existing user who is *not yet verified*, this check might be different.
             // If a user is already fully registered and tries to use "resend OTP" for signup, this is correct.
            return ResponseEntity.badRequest().body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Email đã được sử dụng cho một tài khoản đã hoàn tất đăng ký.", // Clarified message
                    "data", ""
            ));
        }
        // Consider if a different check is needed, e.g., if user is in a "pending verification" state.

        try {
            String newOtp = generateOTP();
            userService.saveOtp(email, newOtp); // Save OTP for the email attempting to register
            emailService.sendOtpEmail(email, newOtp);

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "OTP mới đã được gửi lại vào email của bạn",
                    "data", ""
            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi gửi lại OTP: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    private String generateOTP() {
        return String.valueOf(100000 + new Random().nextInt(900000));
    }
}