package com.ecommerce.command.dto;

public class InventoryResponse {
    private Long id;
    private Integer currentInventory;
    private String status;
    private String message;
    
    // Constructors
    public InventoryResponse() {
    }
    
    public InventoryResponse(Long id, Integer currentInventory, String status, String message) {
        this.id = id;
        this.currentInventory = currentInventory;
        this.status = status;
        this.message = message;
    }
    
    // Static factory method
    public static InventoryResponse updated(Long productId, Integer currentInventory) {
        return new InventoryResponse(productId, currentInventory, "inventory_updated", "Inventory updated successfully");
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public Integer getCurrentInventory() {
        return currentInventory;
    }
    
    public void setCurrentInventory(Integer currentInventory) {
        this.currentInventory = currentInventory;
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
