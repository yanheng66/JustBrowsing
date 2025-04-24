package com.ecommerce.command.dto;

public class TagResponse {
    private Long id;
    private Long tagId;
    private String status;
    private String message;
    
    // Constructors
    public TagResponse() {
    }
    
    public TagResponse(Long id, Long tagId, String status, String message) {
        this.id = id;
        this.tagId = tagId;
        this.status = status;
        this.message = message;
    }
    
    // Static factory methods
    public static TagResponse added(Long productId, Long tagId) {
        return new TagResponse(productId, tagId, "tag_added", "Tag added successfully");
    }
    
    public static TagResponse removed(Long productId, Long tagId) {
        return new TagResponse(productId, tagId, "tag_removed", "Tag removed successfully");
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public Long getTagId() {
        return tagId;
    }
    
    public void setTagId(Long tagId) {
        this.tagId = tagId;
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
