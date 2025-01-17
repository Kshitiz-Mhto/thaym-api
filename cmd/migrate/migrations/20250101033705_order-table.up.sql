CREATE TABLE IF NOT EXISTS orders (
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),             -- Unique identifier for the order
  `userId` CHAR(36) NOT NULL,                         -- Foreign key to associate with the user
  `total` DECIMAL(10, 2) NOT NULL,                    -- Total amount for the order
  `subtotal` DECIMAL(10, 2) NOT NULL,                 -- Subtotal before tax and discounts
  `status` ENUM('pending', 'processing', 'shipped', 'completed', 'cancelled', 'refunded') 
      NOT NULL DEFAULT 'pending',                    -- Order status
  `paymentStatus` ENUM('pending', 'paid', 'refunded') 
      NOT NULL DEFAULT 'pending',                    -- Payment status
  `paymentMethod` VARCHAR(50) NOT NULL,              -- Payment method (e.g., Credit Card, PayPal)
  `address` TEXT NOT NULL,                           -- Shipping address
  `currency` CHAR(3) NOT NULL,                       -- ISO 4217 currency code (e.g., USD, EUR)
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp for order creation
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last updated timestamp

  PRIMARY KEY (id),
  FOREIGN KEY (userId) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
