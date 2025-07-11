package com.example.user_service.security;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.stereotype.Component; // Added for @Value injection

import java.text.ParseException;
import java.util.Collections;
import java.util.Date;

@Component // Make it a component to inject values
public class JwtTokenUtil {

    // NEVER use an empty or weak secret key in production!
    // Generate a strong one, e.g., using a secure random generator.
    // Store it securely, e.g., in environment variables or a secrets manager.
    private static final String SECRET_KEY = "UIY89JSPAdXTF7B8P4MQULxr28UEr4vKE7LDH5pmekBqimsQKHAt5Yf3Vo9U3BAmx9xRJ1AqiTetIjx1oUsErbbA3fGH4xTqxc4rVz7gxeR7h45Zj6mX"; // PLEASE REPLACE THIS!

    @Value("${app.jwt.access-token-duration-ms}") // e.g., 1200000 for 20 minutes
    private long accessTokenDurationMs;


    public String createToken(int id, String role) throws JOSEException {
        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(String.valueOf(id))
                .claim("role", role)
                .claim("id", id) // id claim seems redundant if subject is id, but keeping as per original
                .issueTime(new Date())
                .expirationTime(new Date(System.currentTimeMillis() + accessTokenDurationMs))
                .build();

        JWSHeader header = new JWSHeader(JWSAlgorithm.HS256);
        MACSigner signer = new MACSigner(SECRET_KEY.getBytes());

        SignedJWT signedJWT = new SignedJWT(header, claimsSet);
        signedJWT.sign(signer);

        return signedJWT.serialize();
    }

    public static String getRoleFromToken(String token) { // Keep static if used elsewhere without bean injection
        try {
            SignedJWT signedJWT = SignedJWT.parse(token);
            JWTClaimsSet claims = signedJWT.getJWTClaimsSet();
            return claims.getStringClaim("role");
        } catch (ParseException e) {
            // Log error instead of just printStackTrace
            // e.g., log.error("Error parsing role from token", e);
            e.printStackTrace();
        }
        return null;
    }

    public static Integer getIdFromToken(String token) { // Keep static
        try {
            if (token != null && token.startsWith("Bearer ")) {
                token = token.substring(7);
            }
            if (token == null) return null;

            SignedJWT signedJWT = SignedJWT.parse(token);
            JWTClaimsSet claims = signedJWT.getJWTClaimsSet();
            Object idClaim = claims.getClaim("id");
            if (idClaim instanceof Number) {
                return ((Number) idClaim).intValue();
            } else if (idClaim instanceof String) {
                try {
                    return Integer.parseInt((String) idClaim);
                } catch (NumberFormatException e) {
                    e.printStackTrace(); // log error
                    return null;
                }
            }
            // Consider subject as fallback if "id" claim is missing/malformed
            // return Integer.parseInt(claims.getSubject());
        } catch (ParseException e) {
             e.printStackTrace(); // log error
        }
        return null;
    }

    public static Authentication getAuthentication(String token) { // Keep static
        Integer userId = getIdFromToken(token);
        String role = getRoleFromToken(token);

        if (userId != null && role != null) {
            return new UsernamePasswordAuthenticationToken(
                    userId,
                    null,
                    Collections.singletonList(new SimpleGrantedAuthority(role))
            );
        }
        return null;
    }

    public static boolean verifyToken(String token) { // Keep static
        try {
            if (token != null && token.startsWith("Bearer ")) {
                token = token.substring(7);
            }
             if (token == null) return false;

            SignedJWT signedJWT = SignedJWT.parse(token);
            MACVerifier verifier = new MACVerifier(SECRET_KEY.getBytes());
            
            // Also check expiration
            Date expirationTime = signedJWT.getJWTClaimsSet().getExpirationTime();
            if (expirationTime == null || expirationTime.before(new Date())) {
                // System.out.println("Token expired"); // For debugging
                return false; // Token is expired
            }
            
            return signedJWT.verify(verifier);
        } catch (ParseException | JOSEException e) {
            // System.err.println("Token verification failed: " + e.getMessage()); // For debugging
            // e.printStackTrace(); // log error
        }
        return false;
    }
}