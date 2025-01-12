CREATE TABLE IF NOT EXISTS products (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL,
  `image` VARCHAR(255) NOT NULL,
  `price` DECIMAL(10, 2) NOT NULL,
  `currency` CHAR(3) NOT NULL, -- ISO 4217 currency code (e.g., USD, EUR)
  `quantity` INT UNSIGNED NOT NULL,
  `category` VARCHAR(255) NOT NULL, -- Category for the product
  `tags` JSON DEFAULT NULL, -- List of tags stored as JSON
  `isActive` BOOLEAN NOT NULL DEFAULT TRUE, -- Indicates if the product is active
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
