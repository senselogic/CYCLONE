set @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
set @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
set @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

drop schema if exists `TEST`;

create schema if not exists `TEST` default character set utf8mb4 collate utf8mb4_general_ci;

use `TEST`;

drop table if exists `TEST`.`SECTION`;

create table if not exists `TEST`.`SECTION`(
    `Uuid` VARCHAR(36) NOT NULL,
    `Name` TEXT NULL,
    `Text` TEXT NULL,
    primary key( `Uuid` )
    ) engine = InnoDB;

drop table if exists `TEST`.`ARTICLE`;

create table if not exists `TEST`.`ARTICLE`(
    `Uuid` VARCHAR(36) NOT NULL,
    `SectionUuid` VARCHAR(36) NULL,
    `Title` TEXT NULL,
    `Text` TEXT NULL,
    `DateTime` DATETIME NULL,
    primary key( `Uuid` ),
    index `fk_article_section_1_idx`( `SectionUuid` ASC ),
    constraint `fk_article_section_1`
    foreign key( `SectionUuid` )
    references `TEST`.`SECTION`( `Uuid` )
        on delete set null
        on update no action
    ) engine = InnoDB;

set SQL_MODE=@OLD_SQL_MODE;
set FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
set UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
