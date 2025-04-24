package com.ecommerce.command.dto;

import java.time.ZonedDateTime;

public class ApiResponse<T> {
    private ZonedDateTime timestamp;
    private String status;
    private String message;
    private T data;
    
    // Constructors
    public ApiResponse() {
        this.timestamp = ZonedDateTime.now();
    }
    
    public ApiResponse(String status, String message) {
        this();
        this.status = status;
        this.message = message;
    }
    
    public ApiResponse(String status, String message, T data) {
        this(status, message);
        this.data = data;
    }
    
    // Getters and Setters
    public ZonedDateTime getTimestamp() {
        return timestamp;
    }
    
    public void setTimestamp(ZonedDateTime timestamp) {
        this.timestamp = timestamp;
    }
    
    public String getStatus() {
        return status;
    }
    
    public void setStatus(String status) {
        this.status = status;
    }
    
    public String getMessage() {
        return message;
    }
    
    public void setMessage(String message) {
        this.message = message;
    }
    
    public T getData() {
        return data;
    }
    
    public void setData(T data) {
        this.data = data;
    }
    
    // Builder methods
    public static <T> ApiResponse<T> success(String message) {
        return new ApiResponse<>("success", message);
    }
    
    public static <T> ApiResponse<T> success(String message, T data) {
        return new ApiResponse<>("success", message, data);
    }
    
    public static <T> ApiResponse<T> error(String message) {
        return new ApiResponse<>("error", message);
    }
}
