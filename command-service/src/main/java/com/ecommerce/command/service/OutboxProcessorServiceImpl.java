package com.ecommerce.command.service;

import com.ecommerce.command.domain.OutboxEvent;
import com.ecommerce.command.repository.OutboxRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Async;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.ZonedDateTime;
import java.util.List;

@Service
public class OutboxProcessorServiceImpl implements OutboxProcessorService {
    
    private static final Logger log = LoggerFactory.getLogger(OutboxProcessorServiceImpl.class);
    
    private final OutboxRepository outboxRepository;
    private final KafkaTemplate<String, String> kafkaTemplate;
    
    @Value("${outbox.max-items-per-polling:100}")
    private int maxItemsPerPolling;
    
    public OutboxProcessorServiceImpl(OutboxRepository outboxRepository, KafkaTemplate<String, String> kafkaTemplate) {
        this.outboxRepository = outboxRepository;
        this.kafkaTemplate = kafkaTemplate;
    }
    
    @Override
    @Async("taskExecutor")
    @Scheduled(fixedDelayString = "${outbox.polling.interval.ms:1000}")
    @Transactional
    public void processOutboxEvents() {
        List<OutboxEvent> unprocessedEvents = outboxRepository.findUnprocessedEventsOrderByCreatedAt(maxItemsPerPolling);
        
        if (!unprocessedEvents.isEmpty()) {
            log.info("Processing {} outbox events", unprocessedEvents.size());
            
            for (OutboxEvent event : unprocessedEvents) {
                try {
                    // Determine the topic based on the aggregate type
                    String topic = determineTopicName(event.getAggregateType());
                    String key = event.getAggregateId();
                    
                    // Send the event to Kafka
                    kafkaTemplate.send(topic, key, event.getPayload()).get();
                    
                    // Mark as processed
                    event.setProcessed(true);
                    event.setProcessedAt(ZonedDateTime.now());
                    outboxRepository.save(event);
                    
                    log.debug("Processed outbox event: id={}, type={}", event.getId(), event.getEventType());
                } catch (Exception e) {
                    log.error("Error processing outbox event: id={}", event.getId(), e);
                }
            }
        }
    }
    
    private String determineTopicName(String aggregateType) {
        return switch (aggregateType.toLowerCase()) {
            case "product" -> "products";
            case "inventory" -> "inventory";
            case "order" -> "orders";
            default -> throw new IllegalArgumentException("Unknown aggregate type: " + aggregateType);
        };
    }
}