CREATE TABLE IF NOT EXISTS locations(

	`locationID` CHAR(36) PRIMARY KEY,
	`photoID` CHAR(36) NOT NULL,
	`latitude` DECIMAL(10, 8) NOT NULL,
	`longitude` DECIMAL(11, 8) NOT NULL,
	`createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);