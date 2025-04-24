package com.ecommerce.command.service;

import com.ecommerce.command.dto.InventoryResponse;
import com.ecommerce.command.dto.UpdateInventoryRequest;

public interface InventoryService {
    
    /**
     * Updates inventory for a product
     * @param productId The ID of the product
     * @param request The inventory update request
     * @return The updated inventory response
     */
    InventoryResponse updateInventory(Long productId, UpdateInventoryRequest request);
    
    /**
     * Checks if there is sufficient inventory for a product
     * @param productId The ID of the product
     * @param quantity The quantity to check
     * @return true if there is sufficient inventory, false otherwise
     */
    boolean hasSufficientInventory(Long productId, Integer quantity);
    
    /**
     * Decrements inventory for a product
     * @param productId The ID of the product
     * @param quantity The quantity to decrement
     */
    void decrementInventory(Long productId, Integer quantity);
}
