package com.ecommerce.command.controller;

import com.ecommerce.command.dto.*;
import com.ecommerce.command.service.InventoryService;
import com.ecommerce.command.service.ProductService;
import jakarta.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/products")
public class ProductController {
    
    private static final Logger log = LoggerFactory.getLogger(ProductController.class);
    
    private final ProductService productService;
    private final InventoryService inventoryService;
    
    public ProductController(ProductService productService, InventoryService inventoryService) {
        this.productService = productService;
        this.inventoryService = inventoryService;
    }
    
    @PostMapping
    public ResponseEntity<ProductResponse> createProduct(@Valid @RequestBody CreateProductRequest request) {
        log.info("Received request to create product with SKU: {}", request.getSku());
        ProductResponse response = productService.createProduct(request);
        return new ResponseEntity<>(response, HttpStatus.CREATED);
    }
    
    @PutMapping("/{productId}")
    public ResponseEntity<ProductResponse> updateProduct(
            @PathVariable Long productId,
            @Valid @RequestBody UpdateProductRequest request) {
        log.info("Received request to update product with ID: {}", productId);
        ProductResponse response = productService.updateProduct(productId, request);
        return ResponseEntity.ok(response);
    }
    
    @PostMapping("/{productId}/tags")
    public ResponseEntity<TagResponse> addTagToProduct(
            @PathVariable Long productId,
            @Valid @RequestBody TagDto tagDto) {
        log.info("Received request to add tag to product with ID: {}", productId);
        Long tagId = productService.addTagToProduct(productId, tagDto);
        return ResponseEntity.ok(TagResponse.added(productId, tagId));
    }
    
    @DeleteMapping("/{productId}/tags/{tagId}")
    public ResponseEntity<TagResponse> removeTagFromProduct(
            @PathVariable Long productId,
            @PathVariable Long tagId) {
        log.info("Received request to remove tag from product with ID: {}", productId);
        productService.removeTagFromProduct(productId, tagId);
        return ResponseEntity.ok(TagResponse.removed(productId, tagId));
    }
    
    @PutMapping("/{productId}/inventory")
    public ResponseEntity<InventoryResponse> updateInventory(
            @PathVariable Long productId,
            @Valid @RequestBody UpdateInventoryRequest request) {
        log.info("Received request to update inventory for product with ID: {}", productId);
        InventoryResponse response = inventoryService.updateInventory(productId, request);
        return ResponseEntity.ok(response);
    }
}