CREATE TABLE IF NOT EXISTS fotos (
    `photoID` CHAR(36) PRIMARY KEY,
    `postID` CHAR(36) NOT NULL,
    `location` VARCHAR(255),
    `url_photo` VARCHAR(255) NOT NULL,
    FOREIGN KEY (postID) REFERENCES posts(postID)
);
