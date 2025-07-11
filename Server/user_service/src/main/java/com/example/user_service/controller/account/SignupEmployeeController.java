package com.example.user_service.controller.account;

import com.example.user_service.dto.SignupEmployeeRequest;
import com.example.user_service.model.Role;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
// import com.example.user_service.service.EmailService; // Not used in this controller
import com.example.user_service.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1")
public class SignupEmployeeController {

    private final UserService userService; // Made final
    // private final EmailService emailService; // Not used

    @Autowired
    public SignupEmployeeController(UserService userService) { // Removed EmailService
        this.userService = userService;
        // this.emailService = emailService;
    }

    @PostMapping("/create") // Consider a more descriptive path like /admin/create-employee
    public ResponseEntity<Map<String, Object>> registerEmployee(@RequestHeader("Authorization") String token, @RequestBody SignupEmployeeRequest request) { // ResponseEntity updated
        try {
            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
            User user = userService.getUserById(userId);

            if (user == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy người dùng (admin)", // Clarified
                        "data", ""
                ));
            }

            Optional<Role> role = user.getRoles().stream().findFirst();
            if (role.isPresent() && (role.get().getId() == 1 || "ROLE_ADMIN".equals(role.get().getRoleName())) ) { // Check by name too for robustness
                if (userService.existsByUsername(request.getUsername())) {
                    return ResponseEntity.badRequest().body(Map.of(
                            "code", HttpStatus.BAD_REQUEST.value(),
                            "message", "Email đã được sử dụng",
                            "data", ""
                    ));
                }

                userService.registerEmployee(request);

                return ResponseEntity.ok(Map.of(
                        "code", HttpStatus.OK.value(),
                        "message", "Tạo tài khoản nhân viên thành công",
                        "data", "" // Or return created employee DTO if needed
                ));
            } else {
                return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of(
                        "code", HttpStatus.FORBIDDEN.value(),
                        "message", "Bạn không có quyền thực hiện hành động này",
                        "data", ""
                ));
            }

        } catch (IllegalArgumentException e) { // Catch specific exception from JwtTokenUtil if it throws one for bad token
             return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token không hợp lệ hoặc đã hết hạn: " + e.getMessage(),
                    "data", ""
            ));
        }
        catch (Exception e) {
            // Log the exception e.printStackTrace(); or use a proper logger
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi tạo tài khoản nhân viên: " + e.getMessage(),
                    "data", ""
            ));
        }
    }
}