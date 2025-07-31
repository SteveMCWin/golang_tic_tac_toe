.load ./spellfix

drop table if exists spellfix_users;
create virtual table spellfix_users using spellfix1;

insert into spellfix_users(word) select username from users;
