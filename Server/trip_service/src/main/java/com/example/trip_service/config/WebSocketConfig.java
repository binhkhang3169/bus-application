package com.example.trip_service.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.messaging.simp.config.MessageBrokerRegistry;
// BẮT ĐẦU: CÁC IMPORT CẦN THIẾT
import org.springframework.web.socket.config.annotation.EnableWebSocketMessageBroker;
import org.springframework.web.socket.config.annotation.StompEndpointRegistry;
import org.springframework.web.socket.config.annotation.WebSocketMessageBrokerConfigurer;
// KẾT THÚC: CÁC IMPORT CẦN THIẾT

@Configuration
@EnableWebSocketMessageBroker // Lỗi xảy ra ở đây nếu thiếu import
public class WebSocketConfig implements WebSocketMessageBrokerConfigurer { // Lỗi xảy ra ở đây nếu thiếu import

    @Override
    public void configureMessageBroker(MessageBrokerRegistry config) {
        config.enableSimpleBroker("/topic");
        config.setApplicationDestinationPrefixes("/app");
    }

    @Override
    public void registerStompEndpoints(StompEndpointRegistry registry) {
        registry.addEndpoint("/ws").setAllowedOriginPatterns("*").withSockJS();
    }
}