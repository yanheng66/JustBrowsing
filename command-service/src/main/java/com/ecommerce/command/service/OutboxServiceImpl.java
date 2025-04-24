package com.ecommerce.command.service;

import com.ecommerce.command.domain.OutboxEvent;
import com.ecommerce.command.repository.OutboxRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

@Service
public class OutboxServiceImpl implements OutboxService {
    
    private static final Logger log = LoggerFactory.getLogger(OutboxServiceImpl.class);
    
    private final OutboxRepository outboxRepository;
    
    public OutboxServiceImpl(OutboxRepository outboxRepository) {
        this.outboxRepository = outboxRepository;
    }
    
    @Override
    @Transactional(propagation = Propagation.REQUIRED)
    public OutboxEvent createOutboxEvent(String aggregateType, String aggregateId, String eventType, String payload) {
        log.debug("Creating outbox event: type={}, id={}, event={}", aggregateType, aggregateId, eventType);
        
        OutboxEvent outboxEvent = new OutboxEvent();
        outboxEvent.setAggregateType(aggregateType);
        outboxEvent.setAggregateId(aggregateId);
        outboxEvent.setEventType(eventType);
        outboxEvent.setPayload(payload);
        
        return outboxRepository.save(outboxEvent);
    }
}