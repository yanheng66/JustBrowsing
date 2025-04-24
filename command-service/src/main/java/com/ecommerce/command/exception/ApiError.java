package com.ecommerce.command.exception;

import java.time.ZonedDateTime;
import java.util.ArrayList;
import java.util.List;

public class ApiError {
    private ZonedDateTime timestamp;
    private int status;
    private String error;
    private String message;
    private List<String> details;
    private String path;
    
    public ApiError() {
        this.timestamp = ZonedDateTime.now();
        this.details = new ArrayList<>();
    }
    
    public ApiError(int status, String error, String message, String path) {
        this();
        this.status = status;
        this.error = error;
        this.message = message;
        this.path = path;
    }
    
    // Getters and Setters
    public ZonedDateTime getTimestamp() {
        return timestamp;
    }
    
    public void setTimestamp(ZonedDateTime timestamp) {
        this.timestamp = timestamp;
    }
    
    public int getStatus() {
        return status;
    }
    
    public void setStatus(int status) {
        this.status = status;
    }
    
    public String getError() {
        return error;
    }
    
    public void setError(String error) {
        this.error = error;
    }
    
    public String getMessage() {
        return message;
    }
    
    public void setMessage(String message) {
        this.message = message;
    }
    
    public List<String> getDetails() {
        return details;
    }
    
    public void setDetails(List<String> details) {
        this.details = details;
    }
    
    public void addDetail(String detail) {
        this.details.add(detail);
    }
    
    public String getPath() {
        return path;
    }
    
    public void setPath(String path) {
        this.path = path;
    }
}
