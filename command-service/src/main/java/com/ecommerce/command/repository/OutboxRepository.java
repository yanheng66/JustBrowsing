package com.ecommerce.command.repository;

import com.ecommerce.command.domain.OutboxEvent;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface OutboxRepository extends JpaRepository<OutboxEvent, Long> {
    @Query(value = "SELECT * FROM outbox_events WHERE processed = false ORDER BY created_at ASC LIMIT :limit", nativeQuery = true)
    List<OutboxEvent> findUnprocessedEventsOrderByCreatedAt(@Param("limit") int limit);
    
    List<OutboxEvent> findByAggregateTypeAndAggregateIdAndEventType(String aggregateType, String aggregateId, String eventType);
}
