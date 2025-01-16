CREATE TABLE IF NOT EXISTS orderitems (
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),         -- Unique identifier for the order item
  `orderId` CHAR(36) NOT NULL,                    -- Foreign key linking to the orders table
  `productId` CHAR(36) NOT NULL,                  -- Foreign key linking to the products table
  `productName` VARCHAR(255) NOT NULL,            -- Cached product name
  `quantity` INT UNSIGNED NOT NULL,               -- Quantity ordered, minimum of 1
  `price` DECIMAL(10, 2) NOT NULL,                -- Price per unit
  `totalPrice` DECIMAL(10, 2) NOT NULL,           -- Calculated total price (Quantity * Price)
  `currency` CHAR(3) NOT NULL,                    -- ISO 4217 currency code
  `discount` DECIMAL(10, 2) NOT NULL DEFAULT 0,   -- Discount applied to this item
  `tax` DECIMAL(10, 2) NOT NULL DEFAULT 0,        -- Tax applied to this item
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp for item creation
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last updated timestamp

  PRIMARY KEY (id),
  FOREIGN KEY (orderId) REFERENCES orders(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (productId) REFERENCES products(`productId`) ON DELETE CASCADE ON UPDATE CASCADE
);
