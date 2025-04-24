package com.ecommerce.command.repository;

import com.ecommerce.command.domain.Product;
import com.ecommerce.command.domain.ProductTag;
import com.ecommerce.command.domain.Tag;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface ProductTagRepository extends JpaRepository<ProductTag, Long> {
    Optional<ProductTag> findByProductAndTag(Product product, Tag tag);
    boolean existsByProductAndTag(Product product, Tag tag);
    void deleteByProductAndTag(Product product, Tag tag);
}
