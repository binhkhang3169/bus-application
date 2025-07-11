package com.example.user_service.service;

import com.example.user_service.model.RefreshToken;
import com.example.user_service.model.User;
import com.example.user_service.repository.RefreshTokenRepository;
import com.example.user_service.repository.UserRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.time.temporal.ChronoField;
import java.util.Optional;
import java.util.UUID;

@Service
public class RefreshTokenService {

    private static final Logger logger = LoggerFactory.getLogger(RefreshTokenService.class);

    @Value("${app.jwt.refresh-token.default-duration-ms}")
    private Long oneDayDurationMs;

    // Define the practical maximum Instant for MySQL DATETIME(6)
    // Corresponds to 9999-12-31 23:59:59.999999
    // LocalDateTime.of(year, month, day, hour, minute, second, nanoOfSecond)
    // For .999999 seconds, nanoOfSecond is 999999000.
    // However, to be precise with datetime(6), it's better to ensure it doesn't exceed the database capabilities.
    // A slightly safer approach might be to set nanoseconds to 0 if precision issues are a concern,
    // or use a well-known library constant if available, but constructing it as follows is generally fine.
    // For DATETIME(6), the precision is microseconds.
    public static final Instant PRACTICAL_MAX_MYSQL_DATETIME =
            LocalDateTime.of(9999, 12, 31, 23, 59, 59, 999999000).toInstant(ZoneOffset.UTC);


    @Autowired
    private RefreshTokenRepository refreshTokenRepository;

    @Autowired
    private UserRepository userRepository;

    public Optional<RefreshToken> findByToken(String token) {
        return refreshTokenRepository.findByToken(token);
    }

    @Transactional
    public RefreshToken createRefreshToken(Integer userId, boolean rememberMe) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> {
                    logger.error("Lỗi khi tạo refresh token: Không tìm thấy người dùng với id {}", userId);
                    return new RuntimeException("Lỗi: Không tìm thấy người dùng với id " + userId);
                });

        int deletedCount = refreshTokenRepository.deleteByUser(user);
        if (deletedCount > 0) {
            logger.info("Đã xóa {} refresh token cũ của người dùng id {}", deletedCount, userId);
        }

        RefreshToken refreshToken = new RefreshToken();
        refreshToken.setUser(user);
        refreshToken.setToken(UUID.randomUUID().toString());

        if (rememberMe) {
            // For "Remember Me", set expiry to the maximum supported DATETIME value
            refreshToken.setExpiryDate(PRACTICAL_MAX_MYSQL_DATETIME);
            logger.info("Tạo refresh token (remember me) cho user id {} với expiry {}", userId, PRACTICAL_MAX_MYSQL_DATETIME);
        } else {
            Instant expiry = Instant.now().plusMillis(oneDayDurationMs);
            refreshToken.setExpiryDate(expiry);
            logger.info("Tạo refresh token (no remember me) cho user id {} với expiry {}", userId, expiry);
        }

        return refreshTokenRepository.save(refreshToken);
    }

    public RefreshToken verifyExpiration(RefreshToken token) {
        // Check if the token is set to our practical "never expire" date
        // A direct comparison with PRACTICAL_MAX_MYSQL_DATETIME can be done.
        // Or, more generally, any date very far in the future is treated as non-expiring for practical purposes.
        if (PRACTICAL_MAX_MYSQL_DATETIME.equals(token.getExpiryDate())) {
            logger.debug("Refresh token {} for user {} is effectively non-expiring (set to max practical date).", token.getToken(), token.getUser().getId());
            return token;
        }

        if (token.getExpiryDate().compareTo(Instant.now()) < 0) {
            logger.warn("Refresh token {} for user {} đã hết hạn vào {}. Xóa token.",
                    token.getToken(), token.getUser().getId(), token.getExpiryDate());
            refreshTokenRepository.delete(token);
            throw new RuntimeException("Refresh token đã hết hạn (" + token.getToken() + "). Vui lòng đăng nhập lại.");
        }
        logger.debug("Refresh token {} for user {} còn hạn. Expiry: {}", token.getToken(), token.getUser().getId(), token.getExpiryDate());
        return token;
    }

    @Transactional
    public int deleteByUserId(Integer userId) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> {
                    logger.error("Lỗi khi xóa refresh token: Không tìm thấy người dùng với id {}", userId);
                    return new RuntimeException("Lỗi: Không tìm thấy người dùng với id " + userId);
                });
        logger.info("Xóa refresh token cho user id {}", userId);
        return refreshTokenRepository.deleteByUser(user);
    }

    @Transactional
    public void deleteByToken(String token) {
        if (token == null || token.isEmpty()) {
            logger.warn("Attempted to delete refresh token with null or empty token string.");
            return;
        }
        Optional<RefreshToken> existingToken = refreshTokenRepository.findByToken(token);
        if (existingToken.isPresent()) {
            // refreshTokenRepository.deleteByToken(token); // Assuming this method exists and works by string token
            // If deleteByToken is a custom query method that takes a String:
            refreshTokenRepository.delete(existingToken.get()); // Standard JPA delete
            logger.info("Đã xóa refresh token: {}", token);
        } else {
            logger.info("Không tìm thấy refresh token để xóa: {}", token);
        }
        // The original code had refreshTokenRepository.deleteByToken(token) twice.
        // If `deleteByToken(String token)` in your repository is a custom @Query method for deletion,
        // then one call is sufficient. If it's not, and you rely on findByToken then delete, the above is fine.
        // For clarity, I'll assume `deleteByToken(String token)` is meant to be the primary way if it exists:
        // refreshTokenRepository.deleteByToken(token); // This might be what you intended after the check
        // logger.info("Yêu cầu xóa refresh token: {}. Nếu token tồn tại, nó đã được xóa.", token);
    }
}