CREATE VIEW users_key AS
SELECT
    username,
    email
FROM
    users;

CREATE VIEW user_images AS
SELECT
    u.id AS user_id,
    u.username,
    u.email,
    image.id AS image_record_id,
    image.sent_image_id,
    image.received_image_id
FROM
    users u
    JOIN images image ON u.id = image.user_id;

