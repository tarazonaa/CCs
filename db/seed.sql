-- SAMPLE INSERT WITH ENCRYPTION
INSERT INTO users (username, email, password)
    VALUES ('tarazonaa', pgp_sym_encrypt('someemail@gmail.com', 'SOME_PASSWORD'), pgp_sym_encrypt('123', 'SOME_PASSWORD'));

