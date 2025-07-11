package com.example.user_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;

@Component
public class EmailServiceClient {

    private final RestTemplate restTemplate;
    private final String emailServiceUrl;

    public EmailServiceClient(
            @Value("${email.service.url:http://email_service:8085}") String emailServiceUrl) {
        this.restTemplate = new RestTemplate();
        this.emailServiceUrl = emailServiceUrl;
    }

    public boolean sendEmail(String to, String title, String body, String type) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);

            Map<String, String> requestBody = new HashMap<>();
            requestBody.put("to", to);
            requestBody.put("title", title);
            requestBody.put("body", body);
            requestBody.put("type", type);

            HttpEntity<Map<String, String>> request = new HttpEntity<>(requestBody, headers);
            
            // For direct send
            ResponseEntity<Map> response = restTemplate.postForEntity(
                    emailServiceUrl + "/api/v1/email", 
                    request, 
                    Map.class);
                    
            return response.getStatusCode().is2xxSuccessful();
        } catch (Exception e) {
            // Log the exception
            System.err.println("Failed to send email: " + e.getMessage());
            return false;
        }
    }
    
    public boolean queueEmail(String to, String title, String body, String type) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);

            Map<String, String> requestBody = new HashMap<>();
            requestBody.put("to", to);
            requestBody.put("title", title);
            requestBody.put("body", body);
            requestBody.put("type", type);


            HttpEntity<Map<String, String>> request = new HttpEntity<>(requestBody, headers);
            
            // For queued send
            ResponseEntity<Map> response = restTemplate.postForEntity(
                    emailServiceUrl + "/api/v1/email/queue", 
                    request, 
                    Map.class);
                    
            return response.getStatusCode().is2xxSuccessful();
        } catch (Exception e) {
            // Log the exception
            System.err.println("Failed to queue email: " + e.getMessage());
            return false;
        }
    }
}