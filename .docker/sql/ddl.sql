CREATE TABLE IF NOT EXISTS tokens(
       id BIGSERIAL PRIMARY KEY,
       jwt_token VARCHAR(100) DEFAULT NULL,
       steam_token VARCHAR(100),
       claimed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages(
       id BIGSERIAL PRIMARY KEY,
       sender_name VARCHAR(255),
       content VARCHAR(255),
       created_at TIMESTAMP
);
