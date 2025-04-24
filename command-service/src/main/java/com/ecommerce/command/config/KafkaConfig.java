package com.ecommerce.command.config;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.core.ProducerFactory;

@Configuration
public class KafkaConfig {
    
    @Value("${spring.kafka.producer.properties.schema.registry.url}")
    private String schemaRegistryUrl;
    
    @Bean
    public NewTopic productTopic() {
        return TopicBuilder.name("products")
                .partitions(3)
                .replicas(1)
                .build();
    }
    
    @Bean
    public NewTopic inventoryTopic() {
        return TopicBuilder.name("inventory")
                .partitions(3)
                .replicas(1)
                .build();
    }
    
    @Bean
    public NewTopic orderTopic() {
        return TopicBuilder.name("orders")
                .partitions(3)
                .replicas(1)
                .build();
    }
}
