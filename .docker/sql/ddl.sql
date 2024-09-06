CREATE TABLE IF NOT EXISTS tokens(
       id BIGSERIAL PRIMARY KEY,
       jwt_token VARCHAR(100) DEFAULT NULL,
       steam_token VARCHAR(100),
       claimed_at TIMESTAMP
)
