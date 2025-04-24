package com.ecommerce.command.service;

import com.ecommerce.command.domain.Inventory;
import com.ecommerce.command.domain.Product;
import com.ecommerce.command.domain.ProductTag;
import com.ecommerce.command.domain.Tag;
import com.ecommerce.command.dto.CreateProductRequest;
import com.ecommerce.command.dto.ProductResponse;
import com.ecommerce.command.dto.TagDto;
import com.ecommerce.command.dto.UpdateProductRequest;
import com.ecommerce.command.exception.DuplicateResourceException;
import com.ecommerce.command.exception.ResourceNotFoundException;
import com.ecommerce.command.repository.InventoryRepository;
import com.ecommerce.command.repository.ProductRepository;
import com.ecommerce.command.repository.ProductTagRepository;
import com.ecommerce.command.repository.TagRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;

@Service
public class ProductServiceImpl implements ProductService {
    
    private static final Logger log = LoggerFactory.getLogger(ProductServiceImpl.class);
    
    private final ProductRepository productRepository;
    private final TagRepository tagRepository;
    private final ProductTagRepository productTagRepository;
    private final InventoryRepository inventoryRepository;
    private final OutboxService outboxService;
    private final ObjectMapper objectMapper;
    
    public ProductServiceImpl(
            ProductRepository productRepository,
            TagRepository tagRepository,
            ProductTagRepository productTagRepository,
            InventoryRepository inventoryRepository,
            OutboxService outboxService,
            ObjectMapper objectMapper) {
        this.productRepository = productRepository;
        this.tagRepository = tagRepository;
        this.productTagRepository = productTagRepository;
        this.inventoryRepository = inventoryRepository;
        this.outboxService = outboxService;
        this.objectMapper = objectMapper;
    }
    
    @Override
    @Transactional
    public ProductResponse createProduct(CreateProductRequest request) {
        log.info("Creating new product with SKU: {}", request.getSku());
        
        // Check if product with the same SKU already exists
        if (productRepository.existsBySku(request.getSku())) {
            log.warn("Product with SKU '{}' already exists", request.getSku());
            throw new DuplicateResourceException("Product", "sku", request.getSku());
        }
        
        // Create product
        Product product = new Product();
        product.setSku(request.getSku());
        product.setName(request.getName());
        product.setDescription(request.getDescription());
        product.setPrice(request.getPrice());
        
        // Save product
        product = productRepository.save(product);
        
        // Create inventory
        if (request.getInitialInventory() != null && request.getInitialInventory() > 0) {
            log.debug("Setting initial inventory for product with SKU: {}", request.getSku());
            Inventory inventory = new Inventory();
            inventory.setProduct(product);
            inventory.setQuantity(request.getInitialInventory());
            inventoryRepository.save(inventory);
            
            // Publish inventory created event
            publishInventoryEvent(inventory, "created");
        }
        
        // Add tags if provided
        if (request.getTags() != null && !request.getTags().isEmpty()) {
            log.debug("Adding {} tags to product with SKU: {}", request.getTags().size(), request.getSku());
            for (TagDto tagDto : request.getTags()) {
                addTagInternal(product, tagDto);
            }
        }
        
        // Publish product created event
        publishProductEvent(product, "created");
        
        log.info("Product created successfully with ID: {}", product.getId());
        return ProductResponse.created(product.getId(), product.getSku(), product.getName());
    }
    
    @Override
    @Transactional
    public ProductResponse updateProduct(Long productId, UpdateProductRequest request) {
        log.info("Updating product with ID: {}", productId);
        
        // Find product
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new ResourceNotFoundException("Product", "id", productId));
        
        // Update product
        if (request.getName() != null) {
            product.setName(request.getName());
        }
        
        if (request.getDescription() != null) {
            product.setDescription(request.getDescription());
        }
        
        if (request.getPrice() != null) {
            product.setPrice(request.getPrice());
        }
        
        // Save product
        product = productRepository.save(product);
        
        // Publish product updated event
        publishProductEvent(product, "updated");
        
        log.info("Product updated successfully with ID: {}", product.getId());
        return ProductResponse.updated(product.getId());
    }
    
    @Override
    @Transactional
    public Long addTagToProduct(Long productId, TagDto tagDto) {
        log.info("Adding tag '{}' to product with ID: {}", tagDto.getName(), productId);
        
        // Find product
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new ResourceNotFoundException("Product", "id", productId));
        
        // Add tag
        ProductTag productTag = addTagInternal(product, tagDto);
        
        // Publish product updated event
        publishProductEvent(product, "updated");
        
        log.info("Tag added successfully to product with ID: {}", product.getId());
        return productTag.getTag().getId();
    }
    
    @Override
    @Transactional
    public void removeTagFromProduct(Long productId, Long tagId) {
        log.info("Removing tag with ID: {} from product with ID: {}", tagId, productId);
        
        // Find product
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new ResourceNotFoundException("Product", "id", productId));
        
        // Find tag
        Tag tag = tagRepository.findById(tagId)
                .orElseThrow(() -> new ResourceNotFoundException("Tag", "id", tagId));
        
        // Check if product has the tag
        ProductTag productTag = productTagRepository.findByProductAndTag(product, tag)
                .orElseThrow(() -> new ResourceNotFoundException("Tag", "id", tagId));
        
        // Remove tag from product
        productTagRepository.delete(productTag);
        
        // Publish product updated event
        publishProductEvent(product, "updated");
        
        log.info("Tag removed successfully from product with ID: {}", product.getId());
    }
    
    private ProductTag addTagInternal(Product product, TagDto tagDto) {
        // Find or create tag
        Tag tag = tagRepository.findByName(tagDto.getName())
                .orElseGet(() -> {
                    Tag newTag = new Tag();
                    newTag.setName(tagDto.getName());
                    return tagRepository.save(newTag);
                });
        
        // Check if product already has this tag
        Optional<ProductTag> existingTag = productTagRepository.findByProductAndTag(product, tag);
        if (existingTag.isPresent()) {
            log.warn("Product with ID {} already has tag '{}'", product.getId(), tag.getName());
            throw new DuplicateResourceException("Tag", "name", tag.getName());
        }
        
        // Add tag to product
        ProductTag productTag = new ProductTag();
        productTag.setProduct(product);
        productTag.setTag(tag);
        productTag.setTagValue(tagDto.getValue());
        return productTagRepository.save(productTag);
    }
    
    private void publishProductEvent(Product product, String eventType) {
        try {
            String payload = objectMapper.writeValueAsString(product);
            outboxService.createOutboxEvent("product", product.getId().toString(), eventType, payload);
        } catch (JsonProcessingException e) {
            log.error("Error serializing product for event publishing", e);
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