package com.ecommerce.command.domain;

import jakarta.persistence.*;
import java.time.ZonedDateTime;

@Entity
@Table(name = "inventory")
public class Inventory {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @OneToOne
    @JoinColumn(name = "product_id")
    private Product product;
    
    @Column(nullable = false)
    private Integer quantity = 0;
    
    @Version
    @Column(nullable = false)
    private Integer version = 0;
    
    @Column(name = "last_replenishment_at")
    private ZonedDateTime lastReplenishmentAt;
    
    @Column(name = "updated_at")
    private ZonedDateTime updatedAt;
    
    @PrePersist
    @PreUpdate
    public void preUpdate() {
        updatedAt = ZonedDateTime.now();
    }
    
    // Business logic methods
    public void incrementQuantity(int amount) {
        if (amount < 0) {
            throw new IllegalArgumentException("Increment amount must be positive");
        }
        this.quantity += amount;
        this.lastReplenishmentAt = ZonedDateTime.now();
    }
    
    public void decrementQuantity(int amount) {
        if (amount < 0) {
            throw new IllegalArgumentException("Decrement amount must be positive");
        }
        if (this.quantity < amount) {
            throw new IllegalStateException("Insufficient inventory");
        }
        this.quantity -= amount;
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
    
    public Integer getQuantity() {
        return quantity;
    }
    
    public void setQuantity(Integer quantity) {
        if (quantity < 0) {
            throw new IllegalArgumentException("Quantity cannot be negative");
        }
        this.quantity = quantity;
    }
    
    public Integer getVersion() {
        return version;
    }
    
    public void setVersion(Integer version) {
        this.version = version;
    }
    
    public ZonedDateTime getLastReplenishmentAt() {
        return lastReplenishmentAt;
    }
    
    public void setLastReplenishmentAt(ZonedDateTime lastReplenishmentAt) {
        this.lastReplenishmentAt = lastReplenishmentAt;
    }
    
    public ZonedDateTime getUpdatedAt() {
        return updatedAt;
    }
    
    public void setUpdatedAt(ZonedDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }
}
