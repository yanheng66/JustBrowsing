package com.ecommerce.command.dto;

import jakarta.validation.Valid;
import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;

public class CreateProductRequest {
    @NotBlank(message = "SKU cannot be empty")
    private String sku;
    
    @NotBlank(message = "Name cannot be empty")
    private String name;
    
    private String description;
    
    @NotNull(message = "Price cannot be null")
    @DecimalMin(value = "0.01", message = "Price must be greater than zero")
    private BigDecimal price;
    
    @Valid
    private List<TagDto> tags = new ArrayList<>();
    
    private Integer initialInventory = 0;
    
    // Getters and Setters
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
    
    public String getDescription() {
        return description;
    }
    
    public void setDescription(String description) {
        this.description = description;
    }
    
    public BigDecimal getPrice() {
        return price;
    }
    
    public void setPrice(BigDecimal price) {
        this.price = price;
    }
    
    public List<TagDto> getTags() {
        return tags;
    }
    
    public void setTags(List<TagDto> tags) {
        this.tags = tags;
    }
    
    public Integer getInitialInventory() {
        return initialInventory;
    }
    
    public void setInitialInventory(Integer initialInventory) {
        this.initialInventory = initialInventory;
    }
}
