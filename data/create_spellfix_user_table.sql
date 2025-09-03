.load ../extensions/spellfix

DROP TABLE IF EXISTS spellfix_users;
CREATE virtual TABLE spellfix_users USING spellfix1;

INSERT INTO spellfix_users(word) SELECT username FROM users;
