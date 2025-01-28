CREATE TABLE IF NOT EXISTS users(
  `id` CHAR(36) NOT NULL DEFAULT (UUID()),
  `firstName` VARCHAR(255) NOT NULL,  
  `lastName` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password` CHAR(60) NOT NULL, -- For storing bcrypt hashed passwords
  `isVerified` BOOLEAN NOT NULL DEFAULT FALSE, -- Whether the email is verified
  `role` ENUM('user', 'admin', 'storeowner') NOT NULL DEFAULT 'user', -- Role-based access control
  `isLocked` BOOLEAN NOT NULL DEFAULT FALSE, -- Lock the account if needed
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deletedAt` TIMESTAMP NULL DEFAULT NULL, -- Enable soft deletes
  
  PRIMARY KEY (id),
  UNIQUE KEY (email),
  INDEX (`isVerified`), -- Frequently queried column for verified users
  INDEX (`role`) -- Useful for filtering users by role
);