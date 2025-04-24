package com.ecommerce.command.service;

import com.ecommerce.command.dto.CreateProductRequest;
import com.ecommerce.command.dto.ProductResponse;
import com.ecommerce.command.dto.TagDto;
import com.ecommerce.command.dto.UpdateProductRequest;

public interface ProductService {
    
    /**
     * Creates a new product with the given details
     * @param request The product creation request
     * @return The product creation response
     */
    ProductResponse createProduct(CreateProductRequest request);
    
    /**
     * Updates an existing product with the given details
     * @param productId The ID of the product to update
     * @param request The product update request
     * @return The product update response
     */
    ProductResponse updateProduct(Long productId, UpdateProductRequest request);
    
    /**
     * Adds a tag to a product
     * @param productId The ID of the product
     * @param tagDto The tag details
     * @return The tag ID
     */
    Long addTagToProduct(Long productId, TagDto tagDto);
    
    /**
     * Removes a tag from a product
     * @param productId The ID of the product
     * @param tagId The ID of the tag to remove
     */
    void removeTagFromProduct(Long productId, Long tagId);
}
