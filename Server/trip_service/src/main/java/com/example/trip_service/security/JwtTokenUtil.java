package com.example.trip_service.security;


import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.authority.SimpleGrantedAuthority;

import java.text.ParseException;
import java.util.Collections;
import java.util.Date;


public class JwtTokenUtil {

    private static final String SECRET_KEY = "UIY89JSPAdXTF7B8P4MQULxr28UEr4vKE7LDH5pmekBqimsQKHAt5Yf3Vo9U3BAmx9xRJ1AqiTetIjx1oUsErbbA3fGH4xTqxc4rVz7gxeR7h45Zj6mX"; // Secret key


    public static String createToken(int id, String role) throws JOSEException {

        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(String.valueOf(id))
                .claim("role", role)
                .claim("id", id)
                .issueTime(new Date())
                .expirationTime(new Date(System.currentTimeMillis() + 2 * 60 * 1000))

                .build();


        JWSHeader header = new JWSHeader(JWSAlgorithm.HS256);
        MACSigner signer = new MACSigner(SECRET_KEY.getBytes());


        SignedJWT signedJWT = new SignedJWT(header, claimsSet);
        signedJWT.sign(signer);

        return signedJWT.serialize();
    }


    public static String getRoleFromToken(String token) {
        try {
            SignedJWT signedJWT = SignedJWT.parse(token);
            JWTClaimsSet claims = signedJWT.getJWTClaimsSet();
            return claims.getStringClaim("role");
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return null;
    }

    public static Integer getIdFromToken(String token) {
        try {
            if (token.startsWith("Bearer ")) {
                token = token.substring(7);
            }

            SignedJWT signedJWT = SignedJWT.parse(token);
            JWTClaimsSet claims = signedJWT.getJWTClaimsSet();
            return claims.getClaim("id") != null ? ((Number) claims.getClaim("id")).intValue() : null;
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return null;
    }


    public static Authentication getAuthentication(String token) {
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

    public static boolean verifyToken(String token) {
        try {
            SignedJWT signedJWT = SignedJWT.parse(token);
            MACVerifier verifier = new MACVerifier(SECRET_KEY.getBytes());
            return signedJWT.verify(verifier);
        } catch (ParseException | JOSEException e) {
            e.printStackTrace();
        }
        return false;
    }


}