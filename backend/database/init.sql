CREATE DATABASE IF NOT EXISTS forum_nba CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE forum_nba;

DROP TABLE IF EXISTS topic_reactions;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'ouvert',
    author_id INT NOT NULL,
    category_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_topics_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE topic_reactions (
    user_id INT NOT NULL,
    topic_id INT NOT NULL,
    action_type VARCHAR(10) NOT NULL,
    PRIMARY KEY (user_id, topic_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
);

INSERT INTO categories (id, name) VALUES
(1, 'Boston Celtics'), (2, 'Oklahoma City Thunder'), (3, 'San Antonio Spurs'),
(4, 'Denver Nuggets'), (5, 'Los Angeles Lakers'), (6, 'Houston Rockets'),
(7, 'Minnesota Timberwolves'), (8, 'Portland Trail Blazers'), (9, 'Phoenix Suns'),
(10, 'LA Clippers'), (11, 'Golden State Warriors'), (12, 'New Orleans Pelicans'),
(13, 'Dallas Mavericks'), (14, 'Memphis Grizzlies'), (15, 'Sacramento Kings'),
(16, 'Utah Jazz'), (17, 'Detroit Pistons'), (18, 'New York Knicks'),
(19, 'Cleveland Cavaliers'), (20, 'Toronto Raptors'), (21, 'Atlanta Hawks'),
(22, 'Philadelphia 76ers'), (23, 'Orlando Magic'), (24, 'Charlotte Hornets'),
(25, 'Miami Heat'), (26, 'Milwaukee Bucks'), (27, 'Chicago Bulls'),
(28, 'Brooklyn Nets'), (29, 'Indiana Pacers'), (30, 'Washington Wizards');

INSERT INTO users (id, username, email, password, role) VALUES
(1, 'AdminMathys', 'mathys@nbatalkzone.fr', 'faux_hash_temporaire_1', 'admin'),
(2, 'TheodoreDev', 'theodore@nbatalkzone.fr', 'faux_hash_temporaire_2', 'admin'),
(3, 'FanDesLakers', 'fan@lakers.com', 'faux_hash_temporaire_3', 'user');

INSERT INTO topics (id, title, content, status, author_id, category_id) VALUES
(1, 'La saison incroyable des Celtics 🍀', 'Pensez-vous qu''ils peuvent faire le back-to-back cette année ? L''effectif est monstrueux.', 'ouvert', 1, 1),
(2, 'Wemby est injouable', 'Ce que fait Victor pour les Spurs dépasse l''entendement. Ses contres sont hallucinants.', 'ouvert', 2, 3),
(3, 'LeBron va-t-il prendre sa retraite ?', 'Je pense qu''il fera une dernière tournée d''adieux, mais les Lakers doivent recruter.', 'ouvert', 3, 5),
(4, '🚨 Sujet Modéré : Rumeurs de transfert', 'Ce sujet a été verrouillé par l''administration suite à des débordements.', 'fermé', 1, 18);

INSERT INTO topic_reactions (user_id, topic_id, action_type) VALUES
(2, 1, 'like'),
(3, 1, 'like'),
(1, 2, 'like'),
(3, 2, 'dislike'),
(2, 3, 'like');