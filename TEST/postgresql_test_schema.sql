create table if not exists "TEST"."SECTION"(
    "Uuid" VARCHAR(36) NOT NULL,
    "Name" TEXT NULL,
    "Text" TEXT NULL,
    primary key( "Uuid" )
    ) engine = InnoDB;

create table if not exists "TEST"."ARTICLE"(
    "Uuid" VARCHAR(36) NOT NULL,
    "SectionUuid" VARCHAR(36) NULL,
    "Title" TEXT NULL,
    "Text" TEXT NULL,
    "DateTime" DATETIME NULL,
    primary key( "Uuid" ),
    index `index_article_section_1_idx`( "SectionUuid" ASC )
    ) engine = InnoDB;
