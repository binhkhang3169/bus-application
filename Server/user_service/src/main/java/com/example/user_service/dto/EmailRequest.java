// file: com/example/user_service/dto/EmailRequest.java
package com.example.user_service.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class EmailRequest {
    private String to;
    private String title;
    private String body; // Sẽ chứa chuỗi JSON của dữ liệu cho template
    private String type;
}