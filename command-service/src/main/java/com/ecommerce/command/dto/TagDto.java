package com.ecommerce.command.dto;

import jakarta.validation.constraints.NotBlank;

public class TagDto {
    @NotBlank(message = "Tag name cannot be empty")
    private String name;
    
    @NotBlank(message = "Tag value cannot be empty")
    private String value;
    
    // Constructors
    public TagDto() {
    }
    
    public TagDto(String name, String value) {
        this.name = name;
        this.value = value;
    }
    
    // Getters and Setters
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
    
    public String getValue() {
        return value;
    }
    
    public void setValue(String value) {
        this.value = value;
    }
}
