// file: com/example/user_service/service/EmailService.java
package com.example.user_service.service;

import com.example.user_service.client.EmailKafkaProducer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Map;

@Service
public class EmailService {

    @Autowired
    private EmailKafkaProducer emailKafkaProducer; // << REPLACED: Sử dụng Kafka producer

    /**
     * Gửi email chứa thông tin tài khoản cho nhân viên mới.
     * Cần có template "staff_account_info.html" bên email-service để hoạt động tốt nhất.
     */
    public void sendAccountInfoEmail(String to, String username, String password) {
        String title = "Thông tin tài khoản nhân viên";
        
        // Lấy tên từ email (phần trước @)
        String customerName = extractNameFromEmail(to);
        
        // Lấy ngày hiện tại
        String currentDate = getCurrentFormattedDate();
        
        // Tạo payload cho template
        Map<String, String> data = Map.of(
            "customerName", customerName,
            "username", username,
            "password", password,
            "currentDate", currentDate
        );
        String type = "staff_account_info"; // Loại email mới

        emailKafkaProducer.queueEmailRequest(to, title, data, type);
    }

    /**
     * Gửi email chứa mã OTP.
     */
    public void sendOtpEmail(String to, String otp) {
        String title = "Xác thực tài khoản - Mã OTP của bạn";
        
        // Lấy tên từ email (phần trước @)
        String customerName = extractNameFromEmail(to);
        
        // Lấy ngày hiện tại
        String currentDate = getCurrentFormattedDate();
        
        // email-service (Go) mong đợi một JSON object
        Map<String, String> data = Map.of(
            "data", otp,
            "customerName", customerName,
            "currentDate", currentDate
        );
        String type = "otp";
        
        emailKafkaProducer.queueEmailRequest(to, title, data, type);
    }
    
    /**
     * Gửi một email văn bản đơn giản.
     * Dùng type "generic" và cần template tương ứng bên email-service.
     */
    public void sendEmail(String to, String subject, String text) {
        String type = "generic_text"; // Loại email mới cho văn bản thô
        Map<String, String> data = Map.of("data", text);

        emailKafkaProducer.queueEmailRequest(to, subject, data, type);
    }
    
    /**
     * Helper method để lấy tên từ email (phần trước @)
     */
    private String extractNameFromEmail(String email) {
        if (email != null && email.contains("@")) {
            return email.substring(0, email.indexOf("@"));
        }
        return "User"; // fallback nếu email không hợp lệ
    }
    
    /**
     * Helper method để lấy ngày hiện tại với format đẹp
     */
    private String getCurrentFormattedDate() {
        LocalDate currentDate = LocalDate.now();
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("dd MMM, yyyy");
        return currentDate.format(formatter);
    }
}