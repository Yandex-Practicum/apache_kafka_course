CREATE TABLE users (
    id int PRIMARY KEY,
    name varchar(255),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    private_info VARCHAR
);

TRUNCATE users;

INSERT INTO users (id, name, private_info)
SELECT
    i,
    'Name_' || i || '_' || substring('abcdefghijklmnopqrstuvwxyz', (random() * 26)::integer + 1, 1),
    'Private_info_' || i || '_' || substring('abcdefghijklmnopqrstuvwxyz', (random() * 26)::integer + 1, 1)
FROM
    generate_series(1, 9000000) AS i;

SELECT COUNT(*) FROM users;