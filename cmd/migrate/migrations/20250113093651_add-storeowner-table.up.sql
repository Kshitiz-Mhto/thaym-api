CREATE TABLE IF NOT EXISTS storeowners (
  `storeId` CHAR(36) NOT NULL UNIQUE,         -- Unique store identifier (UUID)
  `name` VARCHAR(255) NOT NULL,              -- Store name
  `ownerName` VARCHAR(255) NOT NULL,         -- Store owner's name
  `email` VARCHAR(255) NOT NULL UNIQUE,      -- Store owner's email
  `phone` VARCHAR(20) DEFAULT NULL,          -- Contact phone number
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Update timestamp
  PRIMARY KEY (`storeId`)                    -- Primary key is storeId
);
