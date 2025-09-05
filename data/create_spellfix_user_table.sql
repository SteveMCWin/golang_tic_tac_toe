.load extensions/spellfix -- has to be extensions/spellfix and not ../extensions/spellfix because of docker

DROP TABLE IF EXISTS spellfix_users;
CREATE virtual TABLE spellfix_users USING spellfix1;

INSERT INTO spellfix_users(word) SELECT username FROM users;
