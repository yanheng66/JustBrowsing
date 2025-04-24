package com.ecommerce.command.service;

import com.ecommerce.command.domain.Order;
import com.ecommerce.command.domain.OrderItem;
import com.ecommerce.command.domain.Product;
import com.ecommerce.command.dto.CreateOrderRequest;
import com.ecommerce.command.dto.OrderItemRequest;
import com.ecommerce.command.dto.OrderResponse;
import com.ecommerce.command.exception.InsufficientInventoryException;
import com.ecommerce.command.exception.ResourceNotFoundException;
import com.ecommerce.command.repository.OrderRepository;
import com.ecommerce.command.repository.ProductRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.concurrent.ThreadLocalRandom;

@Service
public class OrderServiceImpl implements OrderService {
    
    private static final Logger log = LoggerFactory.getLogger(OrderServiceImpl.class);
    
    private final OrderRepository orderRepository;
    private final ProductRepository productRepository;
    private final InventoryService inventoryService;
    private final OutboxService outboxService;
    private final ObjectMapper objectMapper;
    
    public OrderServiceImpl(
            OrderRepository orderRepository,
            ProductRepository productRepository,
            InventoryService inventoryService,
            OutboxService outboxService,
            ObjectMapper objectMapper) {
        this.orderRepository = orderRepository;
        this.productRepository = productRepository;
        this.inventoryService = inventoryService;
        this.outboxService = outboxService;
        this.objectMapper = objectMapper;
    }
    
    @Override
    @Transactional
    public OrderResponse createOrder(CreateOrderRequest request) {
        log.info("Creating new order with {} items", request.getItems().size());
        
        // Generate order number
        String orderNumber = generateOrderNumber();
        
        // Create order
        Order order = new Order();
        order.setOrderNumber(orderNumber);
        order.setTotalAmount(java.math.BigDecimal.ZERO);
        
        // Check and reserve inventory for all items
        for (OrderItemRequest itemRequest : request.getItems()) {
            log.debug("Processing order item: productId={}, quantity={}", itemRequest.getProductId(), itemRequest.getQuantity());
            
            // Check if product exists
            Product product = productRepository.findById(itemRequest.getProductId())
                    .orElseThrow(() -> new ResourceNotFoundException("Product", "id", itemRequest.getProductId()));
            
            // Check if there is sufficient inventory
            if (!inventoryService.hasSufficientInventory(product.getId(), itemRequest.getQuantity())) {
                log.warn("Insufficient inventory for product ID {}", product.getId());
                throw new InsufficientInventoryException(product.getId(), itemRequest.getQuantity(), 
                        product.getInventory() != null ? product.getInventory().getQuantity() : 0);
            }
        }
        
        // Process all items
        for (OrderItemRequest itemRequest : request.getItems()) {
            Product product = productRepository.findById(itemRequest.getProductId()).get(); // Safe as we've already checked
            
            // Create order item
            OrderItem orderItem = new OrderItem();
            orderItem.setProduct(product);
            orderItem.setQuantity(itemRequest.getQuantity());
            orderItem.setUnitPrice(product.getPrice());
            
            // Add item to order
            order.addItem(orderItem);
            
            // Decrement inventory
            inventoryService.decrementInventory(product.getId(), itemRequest.getQuantity());
        }
        
        // Save order
        order = orderRepository.save(order);
        
        // Publish order created event
        publishOrderEvent(order, "created");
        
        log.info("Order created successfully with ID: {}, number: {}", order.getId(), order.getOrderNumber());
        return OrderResponse.created(order.getId(), order.getOrderNumber(), order.getTotalAmount());
    }
    
    @Override
    public String generateOrderNumber() {
        LocalDate today = LocalDate.now();
        String datePart = today.format(DateTimeFormatter.ofPattern("yyyyMMdd"));
        int randomNum = ThreadLocalRandom.current().nextInt(10000, 100000);
        return String.format("ORD-%s-%d", datePart, randomNum);
    }
    
    private void publishOrderEvent(Order order, String eventType) {
        try {
            String payload = objectMapper.writeValueAsString(order);
            outboxService.createOutboxEvent("order", order.getId().toString(), eventType, payload);
        } catch (JsonProcessingException e) {
            log.error("Error serializing order for event publishing", e);
        }
    }
}