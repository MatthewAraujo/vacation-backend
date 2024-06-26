CREATE TABLE IF NOT EXISTS photos( 
    `photoID` CHAR(36) PRIMARY KEY,
    `postID` CHAR(36) NOT NULL,
    `url_photo` TEXT NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postID) REFERENCES posts(postID) ON DELETE CASCADE
)