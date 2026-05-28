INSERT INTO users (email, password_hash)
VALUES
    ('alice@example.com', '$2a$10$OoD26KFtro1xNQ4BDJ26wuGJi0bSR4CxghMeVhVLr.lxdKZsduI0C'),
    ('bob@example.com', '$2a$10$OoD26KFtro1xNQ4BDJ26wuGJi0bSR4CxghMeVhVLr.lxdKZsduI0C')
ON CONFLICT (email) DO UPDATE
SET password_hash = EXCLUDED.password_hash;

INSERT INTO rooms (name)
VALUES
    ('Salle Atlas'),
    ('Salle Neptune'),
    ('Salle Orion'),
    ('Salle Polaris')
ON CONFLICT (name) DO NOTHING;
