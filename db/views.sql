-- SAMPLE VIEW TO SEE USERS
CREATE VIEW users_key AS
SELECT
    username,
    pgp_sym_decrypt(email, "SOME_PASSWORD") AS email,
FROM
    users;

