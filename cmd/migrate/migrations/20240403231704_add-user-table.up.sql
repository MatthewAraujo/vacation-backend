CREATE TABLE IF NOT EXISTS users(
	`id` CHAR(36) UNSIGNED NOT NULL
   	`username` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `createdAt`TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- outras informações do usuário podem ser adicionadas aqui
    UNIQUE (username),
    UNIQUE (email)

	PRIMARY KEY (`id`)
);
