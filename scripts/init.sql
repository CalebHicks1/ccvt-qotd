use qotd;

drop table if exists questions;
CREATE TABLE questions(
    id integer primary key auto_increment,
    body text not null,
    author text default null,
    date_submitted datetime
);

drop table if exists answers;
create table answers (
    id integer primary key auto_increment,
    body text not null,
    author text default null,
    date_submitted datetime,
    votes integer default 0,
    question_id integer not null,
    approved int default 0
);