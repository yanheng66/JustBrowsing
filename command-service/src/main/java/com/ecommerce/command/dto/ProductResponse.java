package com.ecommerce.command.dto;

import java.math.BigDecimal;

public class ProductResponse {
    private Long id;
    private String sku;
    private String name;
    private String status;
    private String message;
    
    // Constructors
    public ProductResponse() {
    }
    
    public ProductResponse(Long id, String sku, String name, String status, String message) {
        this.id = id;
        this.sku = sku;
        this.name = name;
        this.status = status;
        this.message = message;
    }
    
    // Static factory methods
    public static ProductResponse created(Long id, String sku, String name) {
        return new ProductResponse(id, sku, name, "created", "Product created successfully");
    }
    
    public static ProductResponse updated(Long id) {
        return new ProductResponse(id, null, null, "updated", "Product updated successfully");
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public String getSku() {
        return sku;
    }
    
    public void setSku(String sku) {
        this.sku = sku;
    }
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
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
}
