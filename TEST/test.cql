drop keyspace if exists TEST;
create keyspace if not exists TEST with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
drop table if exists TEST.SECTION;
create table if not exists TEST.SECTION( Uuid uuid, Name text, Text text, primary key( Uuid ) );
drop table if exists TEST.ARTICLE;
create table if not exists TEST.ARTICLE( Uuid uuid, SectionUuid uuid, Title text, Text text, DateTime timestamp, primary key( Uuid ) );
create index on TEST.ARTICLE ( SectionUuid );
