package com.ecommerce.command.service;

import com.ecommerce.command.dto.CreateOrderRequest;
import com.ecommerce.command.dto.OrderResponse;

public interface OrderService {
    
    /**
     * Creates a new order with the given items
     * @param request The order creation request
     * @return The order creation response
     */
    OrderResponse createOrder(CreateOrderRequest request);
    
    /**
     * Generates a unique order number
     * @return A unique order number
     */
    String generateOrderNumber();
}
