# API Documentation

This document provides a comprehensive API specification for the CQRS architecture-based e-commerce platform, including all interfaces for the Command Service (write side) and Query Service (read side).

## Command Service API

The Command Service API is responsible for handling all operations that modify the system state, including product management and order processing. All Command APIs return responses in JSON format.

### Product Management

#### Create Product

Creates a new product record.

**Request**

```
POST /api/commands/products
```

**Request Headers**

| Name          | Type   | Description                |
| ------------- | ------ | -------------------------- |
| Content-Type  | string | Must be `application/json` |
| Authorization | string | Bearer token               |

**Request Body**

```json
{
  "sku": "JWLRY-001",
  "name": "Silver Choker Necklace",
  "description": "Elegant silver choker necklace suitable for all occasions",
  "price": 299.99,
  "tags": [
    {
      "name": "brand",
      "value": "Elegant Silver"
    },
    {
      "name": "material",
      "value": "925 silver"
    },
    {
      "name": "color",
      "value": "silver"
    },
    {
      "name": "type",
      "value": "choker necklace"
    }
  ],
  "initialInventory": 100
}
```

**Request Parameters**

| Field            | Type   | Required | Description                                |
| ---------------- | ------ | -------- | ------------------------------------------ |
| sku              | string | Yes      | Product's unique stock keeping unit number |
| name             | string | Yes      | Product name                               |
| description      | string | No       | Product description                        |
| price            | number | Yes      | Product price                              |
| tags             | array  | No       | List of product tags                       |
| tags[].name      | string | Yes      | Tag name                                   |
| tags[].value     | string | Yes      | Tag value                                  |
| initialInventory | number | No       | Initial inventory quantity, defaults to 0  |

**Response**

Success status code: 201 Created

```json
{
  "id": "12345",
  "sku": "JWLRY-001",
  "name": "Silver Choker Necklace",
  "status": "created",
  "message": "Product created successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code      | Description                                            |
| ----------- | --------------- | ------------------------------------------------------ |
| 400         | INVALID_REQUEST | Request format is incorrect or missing required fields |
| 409         | DUPLICATE_SKU   | SKU is already in use                                  |
| 500         | INTERNAL_ERROR  | Server internal error                                  |

#### Update Product

Updates basic information of an existing product.

**Request**

```
PUT /api/commands/products/{productId}
```

**Path Parameters**

| Name      | Type   | Description |
| --------- | ------ | ----------- |
| productId | string | Product ID  |

**Request Headers**

| Name          | Type   | Description                |
| ------------- | ------ | -------------------------- |
| Content-Type  | string | Must be `application/json` |
| Authorization | string | Bearer token               |

**Request Body**

```json
{
  "name": "Premium Silver Choker Necklace",
  "description": "Elegant premium silver choker necklace suitable for formal occasions",
  "price": 329.99
}
```

**Response**

Success status code: 200 OK

```json
{
  "id": "12345",
  "status": "updated",
  "message": "Product updated successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code        | Description                 |
| ----------- | ----------------- | --------------------------- |
| 400         | INVALID_REQUEST   | Request format is incorrect |
| 404         | PRODUCT_NOT_FOUND | Product doesn't exist       |
| 500         | INTERNAL_ERROR    | Server internal error       |

#### Add Product Tag

Adds a new tag to an existing product.

**Request**

```
POST /api/commands/products/{productId}/tags
```

**Path Parameters**

| Name      | Type   | Description |
| --------- | ------ | ----------- |
| productId | string | Product ID  |

**Request Headers**

| Name          | Type   | Description                |
| ------------- | ------ | -------------------------- |
| Content-Type  | string | Must be `application/json` |
| Authorization | string | Bearer token               |

**Request Body**

```json
{
  "name": "collection",
  "value": "Summer New Arrival"
}
```

**Response**

Success status code: 200 OK

```json
{
  "id": "12345",
  "tagId": "789",
  "status": "tag_added",
  "message": "Tag added successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code        | Description                             |
| ----------- | ----------------- | --------------------------------------- |
| 400         | INVALID_REQUEST   | Request format is incorrect             |
| 404         | PRODUCT_NOT_FOUND | Product doesn't exist                   |
| 409         | DUPLICATE_TAG     | The tag already exists for this product |
| 500         | INTERNAL_ERROR    | Server internal error                   |

#### Remove Product Tag

Removes a specific tag from a product.

**Request**

```
DELETE /api/commands/products/{productId}/tags/{tagId}
```

**Path Parameters**

| Name      | Type   | Description |
| --------- | ------ | ----------- |
| productId | string | Product ID  |
| tagId     | string | Tag ID      |

**Request Headers**

| Name          | Type   | Description  |
| ------------- | ------ | ------------ |
| Authorization | string | Bearer token |

**Response**

Success status code: 200 OK

```json
{
  "id": "12345",
  "tagId": "789",
  "status": "tag_removed",
  "message": "Tag removed successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code        | Description                        |
| ----------- | ----------------- | ---------------------------------- |
| 404         | PRODUCT_NOT_FOUND | Product doesn't exist              |
| 404         | TAG_NOT_FOUND     | Tag doesn't exist for this product |
| 500         | INTERNAL_ERROR    | Server internal error              |

#### Update Inventory

Updates product inventory quantity.

**Request**

```
PUT /api/commands/products/{productId}/inventory
```

**Path Parameters**

| Name      | Type   | Description |
| --------- | ------ | ----------- |
| productId | string | Product ID  |

**Request Headers**

| Name          | Type   | Description                |
| ------------- | ------ | -------------------------- |
| Content-Type  | string | Must be `application/json` |
| Authorization | string | Bearer token               |

**Request Body**

```json
{
  "quantityChange": 50,
  "reason": "Restock"
}
```

**Request Parameters**

| Field          | Type   | Required | Description                                                  |
| -------------- | ------ | -------- | ------------------------------------------------------------ |
| quantityChange | number | Yes      | Inventory change quantity, positive to increase, negative to decrease |
| reason         | string | No       | Reason for change                                            |

**Response**

Success status code: 200 OK

```json
{
  "id": "12345",
  "currentInventory": 150,
  "status": "inventory_updated",
  "message": "Inventory updated successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code             | Description                          |
| ----------- | ---------------------- | ------------------------------------ |
| 400         | INVALID_REQUEST        | Request format is incorrect          |
| 404         | PRODUCT_NOT_FOUND      | Product doesn't exist                |
| 409         | INSUFFICIENT_INVENTORY | Insufficient inventory for reduction |
| 500         | INTERNAL_ERROR         | Server internal error                |

### Order Management

#### Create Order

Creates a new order.

**Request**

```
POST /api/commands/orders
```

**Request Headers**

| Name          | Type   | Description                |
| ------------- | ------ | -------------------------- |
| Content-Type  | string | Must be `application/json` |
| Authorization | string | Bearer token               |

**Request Body**

```json
{
  "items": [
    {
      "productId": "12345",
      "quantity": 2
    },
    {
      "productId": "67890",
      "quantity": 1
    }
  ]
}
```

**Request Parameters**

| Field             | Type   | Required | Description                               |
| ----------------- | ------ | -------- | ----------------------------------------- |
| items             | array  | Yes      | List of order items                       |
| items[].productId | string | Yes      | Product ID                                |
| items[].quantity  | number | Yes      | Purchase quantity, must be greater than 0 |

**Response**

Success status code: 201 Created

```json
{
  "orderId": "ORDER-12345",
  "orderNumber": "ORD-20250422-12345",
  "status": "created",
  "totalAmount": 929.97,
  "message": "Order created successfully"
}
```

**Possible Error Responses**

| Status Code | Error Code             | Description                                      |
| ----------- | ---------------------- | ------------------------------------------------ |
| 400         | INVALID_REQUEST        | Request format is incorrect                      |
| 404         | PRODUCT_NOT_FOUND      | One or more products don't exist                 |
| 409         | INSUFFICIENT_INVENTORY | One or more products have insufficient inventory |
| 500         | INTERNAL_ERROR         | Server internal error                            |

## Query Service API

The Query Service API is responsible for handling all read-only operations that do not modify the system state. All Query APIs return responses in JSON format.

### Product Queries

#### Get Product

Retrieves detailed information about a single product.

**Request**

```
GET /api/queries/products/{productId}
```

**Path Parameters**

| Name      | Type   | Description |
| --------- | ------ | ----------- |
| productId | string | Product ID  |

**Response**

Success status code: 200 OK

```json
{
  "id": "12345",
  "sku": "JWLRY-001",
  "name": "Premium Silver Choker Necklace",
  "description": "Elegant premium silver choker necklace suitable for formal occasions",
  "price": 329.99,
  "tags": [
    {
      "id": "101",
      "name": "brand",
      "value": "Elegant Silver"
    },
    {
      "id": "102",
      "name": "material",
      "value": "925 silver"
    },
    {
      "id": "103",
      "name": "color",
      "value": "silver"
    },
    {
      "id": "104",
      "name": "type",
      "value": "choker necklace"
    },
    {
      "id": "105",
      "name": "collection",
      "value": "Summer New Arrival"
    }
  ],
  "currentInventory": 150,
  "images": [
    "https://example.com/images/jwlry001-1.jpg",
    "https://example.com/images/jwlry001-2.jpg"
  ],
  "created": "2025-04-10T08:15:30Z",
  "updated": "2025-04-22T14:05:12Z"
}
```

**Possible Error Responses**

| Status Code | Error Code        | Description           |
| ----------- | ----------------- | --------------------- |
| 404         | PRODUCT_NOT_FOUND | Product doesn't exist |
| 500         | INTERNAL_ERROR    | Server internal error |

#### Search Products by Tags

Searches for products based on specified tags, returning a list of products that match the criteria. The matching rule requires products to include all tags in the query.

**Request**

```
GET /api/queries/products/search?tags={tagName1}:{tagValue1},{tagName2}:{tagValue2}
```

**Query Parameters**

| Name | Type   | Required | Description                                                  |
| ---- | ------ | -------- | ------------------------------------------------------------ |
| tags | string | Yes      | Tag filter conditions, format is `tagName:tagValue`, multiple tags separated by commas |

**Example**

```
GET /api/queries/products/search?tags=material:925 silver,type:choker necklace
```

**Response**

Success status code: 200 OK

```json
{
  "items": [
    {
      "id": "12345",
      "sku": "JWLRY-001",
      "name": "Premium Silver Choker Necklace",
      "price": 329.99,
      "tags": [
        {
          "name": "brand",
          "value": "Elegant Silver"
        },
        {
          "name": "material",
          "value": "925 silver"
        },
        {
          "name": "color",
          "value": "silver"
        },
        {
          "name": "type",
          "value": "choker necklace"
        },
        {
          "name": "collection",
          "value": "Summer New Arrival"
        }
      ],
      "currentInventory": 150,
      "images": ["https://example.com/images/jwlry001-1.jpg"]
    },
    {
      "id": "67890",
      "sku": "JWLRY-002",
      "name": "Rose Gold Choker Necklace",
      "price": 399.99,
      "tags": [
        {
          "name": "brand",
          "value": "Gold Classic Jewelry"
        },
        {
          "name": "material",
          "value": "925 silver"
        },
        {
          "name": "color",
          "value": "rose gold"
        },
        {
          "name": "type",
          "value": "choker necklace"
        }
      ],
      "currentInventory": 75,
      "images": ["https://example.com/images/jwlry002-1.jpg"]
    }
  ],
  "total": 15
}
```

**Possible Error Responses**

| Status Code | Error Code        | Description                         |
| ----------- | ----------------- | ----------------------------------- |
| 400         | INVALID_PARAMETER | Query parameter format is incorrect |
| 500         | INTERNAL_ERROR    | Server internal error               |

------

## Common Error Response Format

All API error responses use the following JSON format:

```json
{
  "timestamp": "2025-04-22T15:30:45Z",
  "status": 400,
  "error": "INVALID_REQUEST",
  "message": "Invalid request parameters",
  "details": ["Price must be greater than zero", "SKU cannot be empty"],
  "path": "/api/commands/products"
}
```

## Important Notes

1. All timestamps use ISO 8601 format and UTC timezone
2. All prices are in USD, represented with decimal points and up to two decimal places
3. Inventory quantities must be non-negative integers
4. Product SKUs must be globally unique
5. The Command Service enforces strict business rule validation, and all update operations use optimistic locking for concurrency control
6. Data provided by the Query Service may experience brief inconsistencies (eventual consistency model)