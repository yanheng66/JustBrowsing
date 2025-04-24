package com.ecommerce.command.exception;

public class InsufficientInventoryException extends RuntimeException {
    
    private Long productId;
    private Integer requested;
    private Integer available;
    
    public InsufficientInventoryException(Long productId, Integer requested, Integer available) {
        super(String.format("Insufficient inventory for product ID %d. Requested: %d, Available: %d", 
                productId, requested, available));
        this.productId = productId;
        this.requested = requested;
        this.available = available;
    }
    
    public Long getProductId() {
        return productId;
    }
    
    public Integer getRequested() {
        return requested;
    }
    
    public Integer getAvailable() {
        return available;
    }
}
