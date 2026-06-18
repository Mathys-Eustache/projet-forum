# 🏀 NBA TalkZone - Documentation Technique

**NBA TalkZone** est une plateforme de discussion web 100% dédiée aux fans de NBA. Elle permet aux passionnés d'échanger sur les 30 franchises de la ligue au sein de salons dédiés, avec un système de compte sécurisé, de création de sujets et de réactions (likes/dislikes).

Ce projet a été développé en adoptant une approche de **développement agile** et une architecture séparée (Backend/Frontend) communicant par API REST.

---

## 🚀 1. Guide d'Installation et de Lancement

Pour faire tourner le projet en local sur votre machine, suivez ces étapes :

### Prérequis

* **Golang** installé sur votre machine.
* **MySQL** (via XAMPP, WAMP, MAMP, ou en natif) en cours d'exécution.

### Étape 1 : Initialisation de la base de données

1. Ouvrez votre gestionnaire MySQL (DBeaver, phpMyAdmin, ou le terminal).
2. Créez une nouvelle base de données et exécutez le script SQL d'initialisation fourni avec le projet (celui contenant la création des tables `users`, `categories`, `topics`, `topic_reactions` et les `INSERT` des franchises).

### Étape 2 : Configuration de l'environnement (Backend)

Dans le dossier `backend/`, aller dans le dossier `database`et copier coller le fichier init.sql dans votre DBeavers

### Étape 3 : Lancement du Backend (API REST)

Ouvrez un premier terminal, naviguez dans le dossier du backend et lancez le serveur :

```bash
cd backend
go run main.go

```

*Le serveur Backend écoutera sur le port `http://localhost:8080`.*

### Étape 4 : Lancement du Frontend (Serveur Web)

Ouvrez un second terminal, naviguez dans le dossier du frontend et lancez le serveur web :

```bash
cd frontend
go run main.go

```

*Le serveur Frontend écoutera sur le port `http://localhost:8081`.*

### Étape 5 : Accès au site

Ouvrez votre navigateur web et rendez-vous à l'adresse : **`http://localhost:8081`**

---

## 🗂️ 2. Organisation & Suivi du Projet

Pour garantir le respect des délais et une collaboration efficace en binôme, nous avons mis en place des outils de gestion et de partage :

* **Gestion des tâches (Trello) :** Suivi de l'avancement des fonctionnalités en méthode agile, répartition des tickets et gestion des priorités.
👉 [Lien vers le Trello du projet](https://trello.com/b/LBtq3I4G/projet-forum-b1)
* **Partage des ressources (Google Drive) :** Centralisation des maquettes, de la documentation, du diaporama de soutenance et des livrables de conception.
👉 [Lien vers le Google Drive](https://drive.google.com/drive/folders/1DvXJTp2PRlAeVRlZsN658zA0zjVYOHrG?usp=drive_link)

---

## 👥 3. Architecture Globale & Rôles

Le projet est scindé en deux blocs autonomes pour assurer une parallélisation totale du travail.

* **Mathys EUSTACHE (Backend, Sécurité & Persistance) :**
* Architecture & Routage (Go) et configuration globale (CORS).
* Modélisation & SQL (Schéma relationnel MySQL, couche *Repositories* avec requêtes préparées anti-injection).
* Sécurité & Métier (Chiffrement SHA-512, middleware de validation JWT).
* Création des *Controllers* (JSON, mapping DTOs, codes HTTP).


* **Théodore NAJMAN (Frontend, Moteur de Thème & UI) :**
* Serveur d'affichage (Go) pour le routage des vues HTML.
* Intégration Graphique HTML5/CSS3.
* Moteur Graphique Dynamique avec système d'habillage via variables CSS selon l'équipe sélectionnée.
* Logique Applicative (Vanilla JS) : Gestion de session `localStorage` et hydratation dynamique du DOM (Fetch API).



---

## ⚙️ 4. LE BACKEND (Golang - Port 8080)

*Suit l'architecture MVC / N-Tiers (Modèles, Controllers, Services, Repositories).*

* **`backend/main.go` (Point d'entrée) :** Initialise la BDD, injecte les dépendances, configure le routeur et gère le **CORS**.
* **`backend/config/db.go` :** Gère l'accès brut à MySQL. Effectue un `db.Ping()` au démarrage.
* **`models/` & `dto/` :** Représentent la structure des tables MySQL et agissent comme des filtres de transfert (Data Transfer Objects) pour ne pas exposer la structure interne.
* **`repositories/` :** La couche SQL pure. Gère le CRUD des sujets, les jointures complexes (`JOIN`), la pagination (`LIMIT`/`OFFSET`), et la recherche textuelle.
* **`services/` :** La couche métier. Gère le hachage en SHA-512, génère le **Token JWT**, valide la conformité des messages et gère les droits d'édition.
* **`controllers/` & `handlers/` :** Réceptionnent le JSON, et interceptent les requêtes (`auth_middleware.go`) pour vérifier la validité du Token JWT avant d'autoriser une action.

---

## 🎨 5. LE FRONTEND (Serveur Web - Port 8081)

*Un client totalement indépendant qui consomme l'API via des requêtes asynchrones.*

* **`frontend/main.go` :** Serveur web secondaire distribuant les fichiers statiques et compilant les vues avec `html/template`.
* **`templates/` :** Vues HTML (Accueil, Conférences, Formulaires d'authentification). `team.html` sert de squelette dynamique pour le forum d'une franchise.
* **`static/css/` :** * `style.css` : Design global et classes critiques d'affichage (ex: `.hidden-element`).
* `teams.css` : Moteur graphique exploitant les variables CSS (`--team-prim`, `--team-sec`) pour recolorer l'interface selon l'ID de l'équipe sélectionnée.


* **`static/js/` :** * `auth.js` : Gère les appels d'authentification et stocke le Token JWT dans le `localStorage`.
* `app.js` : Moteur d'interactivité (Thématisation, vérification de session, et hydratation dynamique de la liste des sujets via l'API).



---

## 🗄️ 6. LA BASE DE DONNÉES (MySQL)

Le système repose sur un schéma relationnel de 4 tables maîtresses :

1. **`users` :** Gère l'authentification et l'identité (pseudos, emails, mots de passe hachés, rôles).
2. **`categories` :** Définit les salons de discussion (les 30 franchises NBA). Sert de pont entre les sujets et le design du frontend.
3. **`topics` :** La table centrale stockant le fil d'actualité (Titre, contenu, date, statut ouvert/fermé). Elle possède des clés étrangères vers `users` (l'auteur) et `categories` (la franchise visée).
4. **`topic_reactions` :** Table d'association reliant un `user_id` et un `topic_id`. Ce système permet de comptabiliser les likes/dislikes tout en bloquant techniquement les doubles votes d'un même utilisateur.