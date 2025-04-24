package com.ecommerce.command.service;

import com.ecommerce.command.domain.OutboxEvent;

public interface OutboxService {
    
    /**
     * Creates a new outbox event
     * @param aggregateType The type of the aggregate (e.g., "product", "order")
     * @param aggregateId The ID of the aggregate
     * @param eventType The type of the event (e.g., "created", "updated")
     * @param payload The payload of the event in JSON format
     * @return The created outbox event
     */
    OutboxEvent createOutboxEvent(String aggregateType, String aggregateId, String eventType, String payload);
}
