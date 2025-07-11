package com.example.user_service.controller.user;

import com.example.user_service.dto.UserResponseDTO;
import com.example.user_service.model.Role;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1/users") // Changed from "/user/api/v1s" for consistency
public class UserRoleController {

    private final UserService userService; // Made final

    @Autowired
    public UserRoleController(UserService userService) {
        this.userService = userService;
    }

    /**
     * Get all users or users by role
     * Only admin can access this endpoint
     * @param token JWT token for authentication
     * @param roleName Role name to filter users (optional)
     * @return List of users
     */
    @GetMapping("/by-role")
    public ResponseEntity<Map<String, Object>> getUsersByRole( // ResponseEntity updated
            @RequestHeader("Authorization") String token,
            @RequestParam(required = false) String roleName) {
        
        try {
            int userId = JwtTokenUtil.getIdFromToken(token.replace("Bearer ", ""));
            User currentUser = userService.getUserById(userId);
            
            if (currentUser == null) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of(
                        "code", HttpStatus.NOT_FOUND.value(),
                        "message", "Không tìm thấy người dùng (admin)", // Clarified
                        "data", ""
                ));
            }
            
            Optional<Role> userRole = currentUser.getRoles().stream().findFirst();
            if (userRole.isEmpty() || !("ROLE_ADMIN".equals(userRole.get().getRoleName()) || userRole.get().getId() == 1) ) { // Check by name or ID
                return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of(
                        "code", HttpStatus.FORBIDDEN.value(),
                        "message", "Bạn không có quyền thực hiện hành động này",
                        "data", ""
                ));
            }
            
            List<UserResponseDTO> users;
            String message;
            if (roleName != null && !roleName.isEmpty()) {
                users = userService.getUsersByRole(roleName);
                message = "Lấy danh sách người dùng theo vai trò '" + roleName + "' thành công";
            } else {
                users = userService.getAllUsers();
                message = "Lấy tất cả danh sách người dùng thành công";
            }
            
            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", message,
                    "data", users
            ));
            
        } catch (IllegalArgumentException e) {
             return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Token không hợp lệ hoặc đã hết hạn: " + e.getMessage(),
                    "data", ""
            ));
        }
        catch (Exception e) {
            // e.printStackTrace(); // Use a proper logger
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi lấy danh sách người dùng: " + e.getMessage(),
                    "data", "" // Ensure data field is present in error responses too
            ));
        }
    }
}