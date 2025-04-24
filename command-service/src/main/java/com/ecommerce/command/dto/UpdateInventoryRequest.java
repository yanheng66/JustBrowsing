package com.ecommerce.command.dto;

import jakarta.validation.constraints.NotNull;

public class UpdateInventoryRequest {
    @NotNull(message = "Quantity change cannot be null")
    private Integer quantityChange;
    
    private String reason;
    
    // Getters and Setters
    public Integer getQuantityChange() {
        return quantityChange;
    }
    
    public void setQuantityChange(Integer quantityChange) {
        this.quantityChange = quantityChange;
    }
    
    public String getReason() {
        return reason;
    }
    
    public void setReason(String reason) {
        this.reason = reason;
    }
}
