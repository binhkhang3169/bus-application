package com.example.user_service.controller.account;

import com.example.user_service.dto.LoginRequest;
import com.example.user_service.model.RefreshToken;
import com.example.user_service.model.User;
import com.example.user_service.security.JwtTokenUtil;
import com.example.user_service.service.RefreshTokenService;
import com.example.user_service.service.UserService;
import jakarta.servlet.http.HttpServletRequest;
// HttpServletResponse might still be needed if you want to set other types of cookies or headers not related to refresh token
// import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.*;

import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1/auth")
public class LoginController {

    private static final Logger logger = LoggerFactory.getLogger(LoginController.class);


    private final UserService userService;
    private final JwtTokenUtil jwtTokenUtil;
    private final RefreshTokenService refreshTokenService;

    @Value("${app.jwt.refresh-token-header-name:X-Refresh-Token}") // Customizable header name
    private String refreshTokenHeaderName;

    // The properties related to cookie names and durations are no longer directly used in this controller
    // for setting cookies, but RefreshTokenService might still use duration logic for DB expiry.
    // @Value("${app.jwt.refresh-token-cookie-name}")
    // private String refreshTokenCookieName;
    // @Value("${app.jwt.refresh-token.default-duration-ms}")
    // private Long oneDayRefreshTokenDurationMs;


    @Autowired
    public LoginController(UserService userService, JwtTokenUtil jwtTokenUtil, RefreshTokenService refreshTokenService) {
        this.userService = userService;
        this.jwtTokenUtil = jwtTokenUtil;
        this.refreshTokenService = refreshTokenService;
    }

    @PostMapping("/login")
    public ResponseEntity<Map<String, Object>> login(
            @RequestBody LoginRequest loginRequest) { // Removed HttpServletResponse

        User user = userService.findByUsername(loginRequest.getUsername());

        if (user == null) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Tài khoản không tồn tại",
                    "data", ""));
        }

        if (user.getActive() != null && user.getActive() == 0) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(Map.of(
                    "code", HttpStatus.FORBIDDEN.value(),
                    "message", "Tài khoản đã bị khóa",
                    "data", ""));
        }

        if (!userService.authenticateUser(loginRequest.getUsername(), loginRequest.getPassword())) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Sai mật khẩu",
                    "data", ""));
        }

        try {
            String role = userService.getUserRole(loginRequest.getUsername());
            int id = user.getId();
            String accessToken = jwtTokenUtil.createToken(id, role);
            boolean rememberMe = loginRequest.isRememberMe();
            RefreshToken refreshToken = refreshTokenService.createRefreshToken(id, rememberMe);

            // Refresh token is now returned in the body for client-side storage
            Map<String, Object> responseData = Map.of(
                    "accessToken", accessToken,
                    "refreshToken", refreshToken.getToken(), // Client needs to store this
                    "user", Map.of(
                            "id", user.getId(),
                            "username", user.getUsername(),
                            "role", role
                    )
            );

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Đăng nhập thành công",
                    "data", responseData));
        } catch (Exception e) {
            logger.error("Lỗi khi tạo token: {}", e.getMessage(), e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi khi tạo token: " + e.getMessage(),
                    "data", ""));
        }
    }

    @PostMapping("/refresh-token")
    public ResponseEntity<Map<String, Object>> refreshToken(
            HttpServletRequest request,
            @RequestHeader(name = "${app.jwt.refresh-token-header-name:X-Refresh-Token}", required = false) String tokenFromHeader,
            @RequestParam(name = "refreshToken", required = false) String tokenFromParam) {

        String requestRefreshToken = null;

        if (StringUtils.hasText(tokenFromHeader)) {
            requestRefreshToken = tokenFromHeader;
        } else if (StringUtils.hasText(tokenFromParam)) {
            requestRefreshToken = tokenFromParam;
        }

        if (!StringUtils.hasText(requestRefreshToken)) {
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of(
                    "code", HttpStatus.BAD_REQUEST.value(),
                    "message", "Refresh token bị thiếu trong header (" + refreshTokenHeaderName + ") hoặc tham số (refreshToken)",
                    "data", ""));
        }

        final String finalRequestRefreshToken = requestRefreshToken;

        try {
            Optional<RefreshToken> optInitialToken = refreshTokenService.findByToken(finalRequestRefreshToken);

            if (optInitialToken.isEmpty()) {
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                        "code", HttpStatus.UNAUTHORIZED.value(),
                        "message", "Refresh token không tồn tại trong cơ sở dữ liệu.",
                        "data", ""));
            }

            RefreshToken verifiedToken = refreshTokenService.verifyExpiration(optInitialToken.get());
            User user = verifiedToken.getUser();

            if (user == null) {
                logger.error("User not found for a valid refresh token: {}", finalRequestRefreshToken);
                return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                        "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                        "message", "Không tìm thấy thông tin người dùng cho refresh token hợp lệ.",
                        "data", ""));
            }

            if (user.getActive() != null && user.getActive() == 0) {
                refreshTokenService.deleteByToken(finalRequestRefreshToken); // Clean up the now invalid token
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                        "code", HttpStatus.UNAUTHORIZED.value(),
                        "message", "Tài khoản người dùng đã bị khóa. Refresh token đã bị vô hiệu hóa.",
                        "data", ""));
            }

            String role = userService.getUserRole(user.getUsername());
            String newAccessToken = jwtTokenUtil.createToken(user.getId(), role);

            Map<String, String> tokenData = Map.of(
                    "accessToken", newAccessToken,
                    "refreshToken", finalRequestRefreshToken // Return the same refresh token
            );

            return ResponseEntity.ok(Map.of(
                    "code", HttpStatus.OK.value(),
                    "message", "Access token làm mới thành công",
                    "data", tokenData
            ));

        } catch (RuntimeException e) {
            logger.warn("Refresh token không hợp lệ hoặc đã hết hạn: {}. Message: {}", finalRequestRefreshToken, e.getMessage());
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of(
                    "code", HttpStatus.UNAUTHORIZED.value(),
                    "message", "Refresh token không hợp lệ hoặc đã hết hạn: " + e.getMessage(),
                    "data", ""));
        } catch (Exception e) {
            logger.error("Lỗi máy chủ nội bộ khi làm mới token: {}", e.getMessage(), e);
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of(
                    "code", HttpStatus.INTERNAL_SERVER_ERROR.value(),
                    "message", "Lỗi máy chủ nội bộ khi làm mới token: " + e.getMessage(),
                    "data", ""));
        }
    }

    @PostMapping("/logout")
    public ResponseEntity<Map<String, Object>> logoutUser(
            HttpServletRequest request,
            @RequestHeader(name = "${app.jwt.refresh-token-header-name:X-Refresh-Token}", required = false) String tokenFromHeader,
            @RequestParam(name = "refreshToken", required = false) String tokenFromParam) {

        String requestRefreshToken = null;

        if (StringUtils.hasText(tokenFromHeader)) {
            requestRefreshToken = tokenFromHeader;
        } else if (StringUtils.hasText(tokenFromParam)) {
            requestRefreshToken = tokenFromParam;
        }

        if (StringUtils.hasText(requestRefreshToken)) {
            try {
                refreshTokenService.deleteByToken(requestRefreshToken);
                logger.info("Refresh token deleted from DB during logout: (token ending with {})",
                    requestRefreshToken.length() > 4 ? requestRefreshToken.substring(requestRefreshToken.length() - 4) : "****");
            } catch (Exception e) {
                logger.warn("Lỗi khi xóa refresh token (ending with {}) khỏi DB trong quá trình đăng xuất: {}. Client should still clear local token.",
                    requestRefreshToken.length() > 4 ? requestRefreshToken.substring(requestRefreshToken.length() - 4) : "****", e.getMessage());
                // Proceed to return success as the main goal is to inform client to logout
            }
        } else {
            logger.info("Không có refresh token trong header hoặc tham số để xóa khỏi DB trong quá trình đăng xuất.");
        }

        // No cookie to clear from the server side in this approach.
        // Client is responsible for clearing its stored tokens.

        return ResponseEntity.ok(Map.of(
                "code", HttpStatus.OK.value(),
                "message", "Đăng xuất thành công. Hãy đảm bảo client đã xóa token.",
                "data", ""));
    }
}