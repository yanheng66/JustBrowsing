package com.ecommerce.command.service;

import com.ecommerce.command.domain.Inventory;
import com.ecommerce.command.domain.Product;
import com.ecommerce.command.dto.InventoryResponse;
import com.ecommerce.command.dto.UpdateInventoryRequest;
import com.ecommerce.command.exception.InsufficientInventoryException;
import com.ecommerce.command.exception.ResourceNotFoundException;
import com.ecommerce.command.repository.InventoryRepository;
import com.ecommerce.command.repository.ProductRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.ZonedDateTime;

@Service
public class InventoryServiceImpl implements InventoryService {
    
    private static final Logger log = LoggerFactory.getLogger(InventoryServiceImpl.class);
    
    private final InventoryRepository inventoryRepository;
    private final ProductRepository productRepository;
    private final OutboxService outboxService;
    private final ObjectMapper objectMapper;
    
    public InventoryServiceImpl(
            InventoryRepository inventoryRepository,
            ProductRepository productRepository,
            OutboxService outboxService,
            ObjectMapper objectMapper) {
        this.inventoryRepository = inventoryRepository;
        this.productRepository = productRepository;
        this.outboxService = outboxService;
        this.objectMapper = objectMapper;
    }
    
    @Override
    @Transactional
    public InventoryResponse updateInventory(Long productId, UpdateInventoryRequest request) {
        log.info("Updating inventory for product with ID {}, change: {}", productId, request.getQuantityChange());
        
        // Find product
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new ResourceNotFoundException("Product", "id", productId));
        
        // Find or create inventory
        Inventory inventory = inventoryRepository.findByProductIdWithLock(productId)
                .orElseGet(() -> {
                    Inventory newInventory = new Inventory();
                    newInventory.setProduct(product);
                    return newInventory;
                });
        
        // Update inventory
        int quantityChange = request.getQuantityChange();
        if (quantityChange > 0) {
            inventory.incrementQuantity(quantityChange);
            inventory.setLastReplenishmentAt(ZonedDateTime.now());
        } else if (quantityChange < 0) {
            try {
                inventory.decrementQuantity(Math.abs(quantityChange));
            } catch (IllegalStateException e) {
                log.warn("Insufficient inventory for product ID {}", productId);
                throw new InsufficientInventoryException(productId, Math.abs(quantityChange), inventory.getQuantity());
            }
        }
        
        // Save inventory
        inventory = inventoryRepository.save(inventory);
        
        // Publish inventory updated event
        publishInventoryEvent(inventory, "updated");
        
        log.info("Inventory updated successfully for product with ID: {}, new quantity: {}", productId, inventory.getQuantity());
        return InventoryResponse.updated(productId, inventory.getQuantity());
    }
    
    @Override
    @Transactional(readOnly = true)
    public boolean hasSufficientInventory(Long productId, Integer quantity) {
        log.debug("Checking inventory for product ID {}, requested quantity: {}", productId, quantity);
        
        return inventoryRepository.findByProductIdWithLock(productId)
                .map(inventory -> inventory.getQuantity() >= quantity)
                .orElse(false);
    }
    
    @Override
    @Transactional
    public void decrementInventory(Long productId, Integer quantity) {
        log.info("Decrementing inventory for product ID {}, quantity: {}", productId, quantity);
        
        Inventory inventory = inventoryRepository.findByProductIdWithLock(productId)
                .orElseThrow(() -> new ResourceNotFoundException("Inventory", "productId", productId));
        
        try {
            inventory.decrementQuantity(quantity);
            inventoryRepository.save(inventory);
            
            // Publish inventory updated event
            publishInventoryEvent(inventory, "updated");
            
            log.info("Inventory decremented successfully for product ID {}, new quantity: {}", productId, inventory.getQuantity());
        } catch (IllegalStateException e) {
            log.warn("Insufficient inventory for product ID {}", productId);
            throw new InsufficientInventoryException(productId, quantity, inventory.getQuantity());
        }
    }
    
    private void publishInventoryEvent(Inventory inventory, String eventType) {
        try {
            String payload = objectMapper.writeValueAsString(inventory);
            outboxService.createOutboxEvent("inventory", inventory.getId().toString(), eventType, payload);
        } catch (JsonProcessingException e) {
            log.error("Error serializing inventory for event publishing", e);
        }
    }
}