package com.ecommerce.command.dto;

import java.math.BigDecimal;

public class OrderResponse {
    private Long orderId;
    private String orderNumber;
    private String status;
    private BigDecimal totalAmount;
    private String message;
    
    // Constructors
    public OrderResponse() {
    }
    
    public OrderResponse(Long orderId, String orderNumber, String status, BigDecimal totalAmount, String message) {
        this.orderId = orderId;
        this.orderNumber = orderNumber;
        this.status = status;
        this.totalAmount = totalAmount;
        this.message = message;
    }
    
    // Static factory method
    public static OrderResponse created(Long orderId, String orderNumber, BigDecimal totalAmount) {
        return new OrderResponse(orderId, orderNumber, "created", totalAmount, "Order created successfully");
    }
    
    // Getters and Setters
    public Long getOrderId() {
        return orderId;
    }
    
    public void setOrderId(Long orderId) {
        this.orderId = orderId;
    }
    
    public String getOrderNumber() {
        return orderNumber;
    }
    
    public void setOrderNumber(String orderNumber) {
        this.orderNumber = orderNumber;
    }
    
    public String getStatus() {
        return status;
    }
    
    public void setStatus(String status) {
        this.status = status;
    }
    
    public BigDecimal getTotalAmount() {
        return totalAmount;
    }
    
    public void setTotalAmount(BigDecimal totalAmount) {
        this.totalAmount = totalAmount;
    }
    
    public String getMessage() {
        return message;
    }
    
    public void setMessage(String message) {
        this.message = message;
    }
}
