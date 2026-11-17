CREATE DATABASE IF NOT EXISTS `forum_nba` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `forum_nba`;

DROP TABLE IF EXISTS `topic_reactions`;
DROP TABLE IF EXISTS `topics`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(50) NOT NULL UNIQUE,
    `email` VARCHAR(255),
    `password` VARCHAR(255) NOT NULL,
    `role` VARCHAR(50) DEFAULT 'user',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE `categories` (
    `id` INT PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE `topics` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(255) NOT NULL,
    `content` TEXT NOT NULL,
    `status` VARCHAR(50) DEFAULT 'ouvert',
    `likes` INT DEFAULT 0,
    `dislikes` INT DEFAULT 0,
    `user_id` INT NOT NULL,
    `category_id` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `topic_reactions` (
    `user_id` INT NOT NULL,
    `topic_id` INT NOT NULL,
    `reaction_type` VARCHAR(50) NOT NULL,
    PRIMARY KEY (`user_id`, `topic_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`topic_id`) REFERENCES `topics`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`) VALUES
(1, 'test1', 'test1@nba.fr', '$2a$10$Z3xBcEefgHijkLmnopQrstUvWxyZaBcDeFgHiJkLmNoPqRsTuVwXy', 'user'),
(2, 'test2', 'test2@nba.fr', '$2a$10$Z3xBcEefgHijkLmnopQrstUvWxyZaBcDeFgHiJkLmNoPqRsTuVwXy', 'user'),
(3, 'MathysJury', 'mathys@nba.fr', '$2a$10$Z3xBcEefgHijkLmnopQrstUvWxyZaBcDeFgHiJkLmNoPqRsTuVwXy', 'user');

INSERT INTO `categories` (`id`, `name`) VALUES
(1, 'Boston Celtics'),
(2, 'Oklahoma City Thunder'),
(3, 'San Antonio Spurs'),
(4, 'Denver Nuggets'),
(5, 'Los Angeles Lakers'),
(6, 'Houston Rockets'),
(7, 'Minnesota Timberwolves'),
(8, 'Portland Trail Blazers'),
(9, 'Phoenix Suns'),
(10, 'LA Clippers'),
(11, 'Golden State Warriors'),
(12, 'New Orleans Pelicans'),
(13, 'Dallas Mavericks'),
(14, 'Memphis Grizzlies'),
(15, 'Sacramento Kings'),
(16, 'Utah Jazz'),
(17, 'Detroit Pistons'),
(18, 'New York Knicks'),
(19, 'Cleveland Cavaliers'),
(20, 'Toronto Raptors'),
(21, 'Atlanta Hawks'),
(22, 'Philadelphia 76ers'),
(23, 'Orlando Magic'),
(24, 'Charlotte Hornets'),
(25, 'Miami Heat'),
(26, 'Milwaukee Bucks'),
(27, 'Chicago Bulls'),
(28, 'Brooklyn Nets'),
(29, 'Indiana Pacers'),
(30, 'Washington Wizards');

INSERT INTO `topics` (`title`, `content`, `status`, `likes`, `dislikes`, `user_id`, `category_id`) VALUES
('Le 18ème titre est là !', 'Quelle saison incroyable des Celtics, la domination était totale du début à la fin des playoffs.', 'ouvert', 15, 1, 1, 1),
('Tatum vs Brown', 'Qui est le véritable MVP de cette équipe selon vous ? J’ai une préférence pour la régularité de Jaylen.', 'ouvert', 8, 2, 3, 1),
('L’avenir s’annonce radieux', 'SGA est clairement un candidat MVP légitime. Cette équipe est tellement jeune et déjà si forte.', 'ouvert', 12, 0, 2, 2),
('Wemby Saison 2', 'Quelles sont vos attentes pour Victor cette année ? Le cap des 25 points par match est-il atteignable ?', 'ouvert', 24, 0, 1, 3),
('LeBron et Bronny', 'Est-ce que l’association père-fils va fonctionner ou est-ce une simple opération de communication ?', 'ouvert', 5, 9, 3, 5),
('Le successeur de JJ Redick', 'Nouveau coach, nouveaux systèmes. Hâte de voir ce que ça va donner sur le terrain.', 'ouvert', 4, 1, 2, 5),
('La fin d’une ère ?', 'Sans Klay, l’équipe a un visage différent. Heureusement que Steph Curry est toujours là pour envoyer des bombes.', 'ouvert', 19, 3, 1, 11),
('L’ambiance au Madison Square Garden', 'Les playoffs au MSG c’est vraiment autre chose. L’équipe a un cœur énorme cette année.', 'ouvert', 14, 0, 3, 18);