CREATE TABLE IF NOT EXISTS posts( 
    `postID` CHAR(36) PRIMARY KEY,
    `userID` CHAR(36) NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `description` TEXT NOT NULL,
    `favorite` INT, 
    FOREIGN KEY (userID) REFERENCES users(userID)
);