package com.example.user_service.controller.employee;

import com.example.user_service.dto.UserUpdateRequest;
import com.example.user_service.model.Customer;
import com.example.user_service.model.Employee;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.CustomerService;
import com.example.user_service.service.EmployeeService;
import com.example.user_service.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1") // Base path from original file
public class UserEmployeeController {

    private final UserService userService; // Made final
    private final EmployeeService employeeService; // Made final

    @Autowired
    public UserEmployeeController(UserService userService, EmployeeService employeeService) {
        this.userService = userService;
        this.employeeService = employeeService;
    }

    // Changed mapping to be more specific, e.g., /customer/info
    @GetMapping("/employee/info") // Renamed from "/api/v1-info" for clarity
    public ResponseEntity<Map<String, Object>> getUserInfo(@RequestHeader("Authorization") String token) { // ResponseEntity updated
        try {
            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
            User user = userService.getUserById(userId);
            System.out.println(user);
            if (user == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy người dùng",
                        "data", ""
                ));
            }
            Employee employee = employeeService.findById(user.getId()); // Assuming customer ID is same as user ID
            if (employee == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy thông tin khách hàng",
                        "data", ""
                ));
            }


            Map<String, Object> userInfoData = new HashMap<>();
            userInfoData.put("id", user.getId());
            userInfoData.put("username", user.getUsername()); // email
            userInfoData.put("phoneNumber", employee.getPhoneNumber());
            userInfoData.put("fullName", employee.getFullName());
            userInfoData.put("active", user.getActive());
            userInfoData.put("address", employee.getAddress());
            userInfoData.put("gender", employee.getGender());
            userInfoData.put("image", employee.getAvatarUrl());



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

    @PostMapping("/employee/change-image")
    public ResponseEntity<?> changeImage(@RequestHeader("Authorization") String token, @RequestParam String image) {
        try {
            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
            User user = userService.getUserById(userId);
            if (user == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", 404,
                        "message", "Không tìm thấy người dùng"
                ));
            }
            Employee employee = employeeService.findById(user.getId());

            employee.setAvatarUrl(image);
            userService.updateUser(user);

            System.out.println("📸 Ảnh nhận được: " + image);


            return ResponseEntity.ok(Map.of(
                    "code", 200,
                    "message", "Lưu ảnh người dùng thành công!"

            ));
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", 500,
                    "message", "Lỗi khi lấy thông tin người dùng"
            ));
        }
    }

    // Changed mapping to be more specific, e.g., /customer/change-info
    @PostMapping("/employee/change-info") // Renamed from "/change-info"
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
            Employee employee = employeeService.findById(user.getId());
            if (employee == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy thông tin khách hàng để cập nhật",
                        "data", ""
                ));
            }


            employee.setPhoneNumber(updateRequest.getPhoneNumber());
            employee.setAddress(updateRequest.getAddress());
            employee.setGender(updateRequest.getGender());
            employee.setFullName(updateRequest.getFullName());
            // Assuming user.setUsername() or other User fields are not changed here. If they are, call userService.save(user) too.
            employeeService.updateUser(employee); // Assuming this saves the changes

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