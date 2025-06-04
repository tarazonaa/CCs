-- Enable pgcrypto for UUID + encryption
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Insert a test user with encrypted email and password
INSERT INTO users (id, username, email, name,  password)
    VALUES (gen_random_uuid (), 'tarazonaa', 'andres.tara.so@gmail.com', 'Mr. BBT','$2a$12$RhlrMmyvcM0En8PeNOINJus0lE3WKIaRlVD/BTfThF0pqpVYUpKZm')
ON CONFLICT (username)
    DO NOTHING;

-- Insert sample image records for the user
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

-- Insert a consumer
INSERT INTO consumers (id, username, custom_id)
    VALUES (gen_random_uuid (), 'ccs-consumer', 'ccs-id')
ON CONFLICT (custom_id)
    DO NOTHING;

-- Insert an oauth2 client app linked to the consumer
INSERT INTO oauth2_credentials (id, name, client_id, client_secret, redirect_uris, consumer_id)
    VALUES (gen_random_uuid (), 'CCs', 'CCs-client-id', 'holajorge', ARRAY['http://localhost:3000/callback'], (
            SELECT
                id
            FROM
                consumers
            WHERE
                custom_id = 'ccs-id'))
ON CONFLICT (client_id)
    DO NOTHING;

