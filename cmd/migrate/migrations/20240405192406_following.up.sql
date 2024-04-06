CREATE TABLE follow (
    `followID` CHAR(36) PRIMARY KEY,
    `userID` CHAR(36) NOT NULL,
    `userID_following` CHAR(36) NOT NULL,
    FOREIGN KEY (userID) REFERENCES users(userID),
    FOREIGN KEY (userID_following) REFERENCES users(userID)
);
