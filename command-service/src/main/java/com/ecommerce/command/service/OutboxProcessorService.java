package com.ecommerce.command.service;

public interface OutboxProcessorService {
    
    /**
     * Processes unprocessed outbox events by publishing them to Kafka
     */
    void processOutboxEvents();
}
