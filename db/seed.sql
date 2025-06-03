-- SAMPLE INSERT WITH ENCRYPTION
INSERT INTO users (username, email, password)
    VALUES ('tarazonaa', pgp_sym_encrypt('someemail@gmail.com', 'SOME_PASSWORD'), pgp_sym_encrypt('123', 'SOME_PASSWORD'));

INSERT INTO images (user_id, sent_image_id, received_image_id)
    VALUES ((
            SELECT
                id
            FROM
                users
            WHERE
                username = 'tarazonaa'), gen_random_uuid (), gen_random_uuid ()),
    ((
        SELECT
            id
        FROM users
        WHERE
            username = 'tarazonaa'), gen_random_uuid (), gen_random_uuid ());

