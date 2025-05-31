-- SAMPLE VIEW TO SEE USERS
CREATE VIEW users_key AS
SELECT
    username,
    pgp_sym_decrypt(email, "SOME_SALT") AS email,
FROM
    users;

