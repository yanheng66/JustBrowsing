package com.ecommerce.command.domain;

import jakarta.persistence.*;
import java.time.ZonedDateTime;

@Entity
@Table(name = "product_tags", uniqueConstraints = {
    @UniqueConstraint(columnNames = {"product_id", "tag_id"})
})
public class ProductTag {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @ManyToOne
    @JoinColumn(name = "product_id")
    private Product product;
    
    @ManyToOne
    @JoinColumn(name = "tag_id")
    private Tag tag;
    
    @Column(name = "tag_value", length = 255)
    private String tagValue;
    
    @Column(name = "created_at")
    private ZonedDateTime createdAt;
    
    @PrePersist
    public void prePersist() {
        createdAt = ZonedDateTime.now();
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public Product getProduct() {
        return product;
    }
    
    public void setProduct(Product product) {
        this.product = product;
    }
    
    public Tag getTag() {
        return tag;
    }
    
    public void setTag(Tag tag) {
        this.tag = tag;
    }
    
    public String getTagValue() {
        return tagValue;
    }
    
    public void setTagValue(String tagValue) {
        this.tagValue = tagValue;
    }
    
    public ZonedDateTime getCreatedAt() {
        return createdAt;
    }
    
    public void setCreatedAt(ZonedDateTime createdAt) {
        this.createdAt = createdAt;
    }
    
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        
        ProductTag that = (ProductTag) o;
        
        if (product != null ? !product.equals(that.product) : that.product != null) return false;
        return tag != null ? tag.equals(that.tag) : that.tag == null;
    }
    
    @Override
    public int hashCode() {
        int result = product != null ? product.hashCode() : 0;
        result = 31 * result + (tag != null ? tag.hashCode() : 0);
        return result;
    }
}
