
CREATE SCHEMA IF NOT EXISTS `ad_space_auction`;
USE `ad_space_auction`;

-- -----------------------------------------------------
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `space` (
  `id` VARCHAR(45) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `base_price` FLOAT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = 'utf8';

-- -----------------------------------------------------
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `bidder` (
  `id` VARCHAR(45) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = 'utf8';

-- -----------------------------------------------------
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `auction` (
  `id` VARCHAR(45) NOT NULL,
  `space_id` VARCHAR(45) NOT NULL,
  `starts_on` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ends_on` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = 'utf8';

-- -----------------------------------------------------
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `auction_bid` (
  `id` VARCHAR(45) NOT NULL,
  `auction_id` VARCHAR(45) NOT NULL,
  `bidder_id` VARCHAR(45) NOT NULL,
  `bid_price` FLOAT NOT NULL DEFAULT 1,
  `bid_on` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = 'utf8';