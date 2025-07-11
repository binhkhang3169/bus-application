package com.example.user_service.controller.customer;

import com.example.user_service.dto.UserUpdateRequest;
import com.example.user_service.model.Customer;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.CustomerService;
import com.example.user_service.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1") // Base path from original file
public class UserCustomerController {

    private final UserService userService; // Made final
    private final CustomerService customerService; // Made final

    @Autowired
    public UserCustomerController(UserService userService, CustomerService customerService) {
        this.userService = userService;
        this.customerService = customerService;
    }

    // Changed mapping to be more specific, e.g., /customer/info
    @GetMapping("/customer/info") // Renamed from "/api/v1-info" for clarity
    public ResponseEntity<Map<String, Object>> getUserInfo(@RequestHeader("Authorization") String token) { // ResponseEntity updated
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
            Customer customer = customerService.findById(user.getId()); // Assuming customer ID is same as user ID
            if (customer == null) {
                 return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy thông tin khách hàng",
                        "data", ""
                ));
            }


            Map<String, Object> userInfoData = new HashMap<>();
            userInfoData.put("id", user.getId());
            userInfoData.put("username", user.getUsername()); // email
            userInfoData.put("phoneNumber", customer.getPhoneNumber());
            userInfoData.put("fullName", customer.getFullName());
            userInfoData.put("active", user.getActive());
            userInfoData.put("address", customer.getAddress());
            userInfoData.put("gender", customer.getGender());

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Lấy thông tin người dùng thành công",
                    "data", userInfoData
            ));

        } catch (IllegalArgumentException e) { // Catch specific exception from JwtTokenUtil
             return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token không hợp lệ hoặc đã hết hạn: " + e.getMessage(),
                    "data", ""
            ));
        }
        catch (Exception e) {
            // Log the exception e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi lấy thông tin người dùng: " + e.getMessage(),
                    "data", ""
            ));
        }
    }

    // Changed mapping to be more specific, e.g., /customer/change-info
    @PostMapping("/customer/change-info") // Renamed from "/change-info"
    public ResponseEntity<Map<String, Object>> changeInfo(@RequestHeader("Authorization") String token, @RequestBody UserUpdateRequest updateRequest) { // ResponseEntity updated
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
            Customer customer = customerService.findById(user.getId());
            if (customer == null) {
                 return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy thông tin khách hàng để cập nhật",
                        "data", ""
                ));
            }


            customer.setPhoneNumber(updateRequest.getPhoneNumber());
            customer.setAddress(updateRequest.getAddress());
            customer.setGender(updateRequest.getGender());
            customer.setFullName(updateRequest.getFullName());
            // Assuming user.setUsername() or other User fields are not changed here. If they are, call userService.save(user) too.
            customerService.updateUser(customer); // Assuming this saves the changes

            return ResponseEntity.ok(Map.of( // Changed from status(HttpStatus.OK) for consistency
                    "code", HttpStatus.OK.value(),
                    "message", "Lưu thông tin người dùng thành công!",
                    "data", "" // Or return updated customer DTO
            ));
        } catch (IllegalArgumentException e) {
             return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token không hợp lệ hoặc đã hết hạn: " + e.getMessage(),
                    "data", ""
            ));
        }
        catch (Exception e) {
            // Log the exception e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi cập nhật thông tin người dùng: " + e.getMessage(), // Corrected message key
                    "data", ""
            ));
        }
    }
}