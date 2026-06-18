# 🏀 Guide de Soutenance & Documentation Technique - NBA TalkZone

## 👥 1. Organisation Méthodologique & Répartition des Rôles

Le projet a été mené sur un cycle court (une semaine) en adoptant une approche de **développement agile par fonctionnalités**. Pour assurer une parallélisation totale du travail, l'architecture a été scindée en deux blocs autonomes (Backend et Frontend) communicant par API REST.

### 🛠️ Mathys : Spécialisation Backend, Sécurité & Persistance (Port 8080)
* **Architecture & Routage (Go) :** Point d'entrée unique `main.go`, gestion du routeur et configuration globale des en-têtes CORS.
* **Modélisation & SQL :** Conception du schéma relationnel MySQL. Écriture de la couche *Repositories* (requêtes préparées anti-injection, jointures, pagination via `LIMIT`/`OFFSET`).
* **Sécurité & Métier :** Chiffrement des mots de passe en SHA-512. Implémentation du middleware `auth_middleware.go` pour valider cryptographiquement les jetons JWT.
* **Couche d'Échange :** Création des *Controllers* pour réceptionner le JSON, le mapper sur les DTOs, et distribuer les codes HTTP (200, 201, 401).

### 🎨 Théodore : Spécialisation Frontend, Moteur de Thème & Rendu Dynamique (Port 8081)
* **Serveur d'Affichage (Go) :** Serveur web secondaire pour le routage des vues HTML et la compilation des templates (`renderTemplate`).
* **Intégration Graphique :** Design HTML5/CSS3, structure sémantique et responsive.
* **Moteur Graphique Dynamique (`teams.css`) :** Système d'habillage optimisé via variables CSS, permettant de changer toute la charte graphique selon l'ID de l'équipe.
* **Logique Applicative (Vanilla JS) :** Fichiers `auth.js` et `app.js` pour dynamiser l'interface. Gestion de session (`localStorage`), lecture d'URL, et hydratation dynamique du DOM (AJAX/Fetch) avec les *Template Literals*.

---

## ⚙️ 2. LE BACKEND (Golang - Port 8080)
*Suit l'architecture MVC / N-Tiers (Modèle, Vue/Controller, Services, Repositories).*

### 📁 Racine & Configuration
* **`backend/main.go` (Point d'entrée) :**
  Initialise la base de données, injecte les dépendances, configure le routeur et gère le **CORS** pour autoriser le port 8081 à lui parler.
* **`backend/config/db.go` :**
  Gère l'accès brut à MySQL. Effectue un `db.Ping()` au démarrage pour s'assurer que la base répond.

### 📁 Les Données (Modèles & Transfert)
* **`models/` (ex: `user.go`, `topic.go`) :**
  Représentent la structure exacte des tables MySQL en code Go avec les tags JSON pour la sérialisation.
* **`dto/` (Data Transfer Objects) :**
  Agissent comme des filtres. Ils moulent les requêtes entrantes (ex: on ne lit que l'email, le pseudo et le mot de passe lors de l'inscription) pour ne pas exposer directement la base de données.

### 📁 La Couche d'Accès aux Données (Repositories)
*Ne contient que du SQL, aucune logique métier.*
* **`user_repository.go` :** Exécute les requêtes `INSERT` (inscription) et `SELECT` (connexion).
* **`category_repository.go` :** Récupère la liste dynamique des franchises NBA.
* **`topic_repository.go` :** Le cœur du forum. Gère le CRUD des sujets, les jointures complexes (`JOIN users`), la pagination, et la recherche textuelle textuelle (`WHERE title LIKE ?`).

### 📁 La Couche Métier (Services)
*Le "Cerveau" de l'application.*
* **`auth_service.go` :** Gère le hachage en SHA-512 des mots de passe et génère le **Token JWT** signé lors de la connexion.
* **`topic_service.go` :** Valide la conformité des messages avant insertion et gère les droits (vérifie si l'utilisateur a le statut admin/modérateur pour forcer la fermeture d'un topic).

### 📁 La Couche Transport (Controllers & Middlewares)
* **`controllers/` :** Réceptionnent le JSON du Frontend, décodent les données, appellent le bon service, et renvoient les statuts HTTP adéquats (200, 201, 401).
* **`handlers/auth_middleware.go` :** L'intercepteur de sécurité. Extrait le Token JWT du header `Authorization: Bearer`, vérifie sa signature et sa date d'expiration avant de laisser passer la requête.

---

## 🎨 3. LE FRONTEND (Serveur Web - Port 8081)
*Un client totalement indépendant qui "consomme" l'API.*

### 📁 Serveur d'affichage
* **`frontend/main.go` :** Un serveur web secondaire en Go qui distribue les fichiers statiques (CSS/JS) et compile les vues avec `html/template`.

### 📁 Les Vues (Dossier `templates/`)
* **`index.html` :** Page d'accueil globale de NBA TalkZone.
* **`conference-est.html` & `conference-ouest.html` :** Affichent les tableaux des franchises. Les lignes redirigent vers `/team?id=X`.
* **`team.html` :** Vue dynamique principale. Ce n'est qu'un squelette (conteneur) que le JavaScript viendra remplir avec les sujets de l'équipe sélectionnée.
* **`inscription.html` & `connexion.html` :** Formulaires d'authentification.

### 📁 Le Design (Dossier `static/css/`)
* **`style.css` :** Feuille de style globale. Contient les classes critiques `.hidden-element` (`display: none !important;`) essentielles pour permettre au JS de forcer l'affichage ou le masquage des éléments de l'UI.
* **`teams.css` :** Moteur graphique. Exploite les variables CSS (`--team-prim`, `--team-sec`). La classe ajoutée dynamiquement (ex: `.theme-celtics`) surcharge ces variables et recolore instantanément tout le site.

### 📁 La Logique d'Interface (Dossier `static/js/`)
* **`auth.js` :** Bloque le rechargement (`e.preventDefault()`), envoie les données de connexion en POST, et stocke le Token JWT renvoyé par l'API dans le `localStorage` du navigateur.
* **`app.js` :** Le moteur d'interactivité :
  1. *Thématisation :* Lit l'ID dans l'URL et applique la bonne classe CSS au `<body>`.
  2. *Sécurité UI :* Vérifie la présence du token JWT et affiche/cache la zone de saisie.
  3. *Hydratation :* Utilise `fetch()` pour interroger l'API avec les paramètres de tri/recherche, et rebâtit le HTML dynamiquement avec les *Template Literals*.

---

## 🗄️ 4. LA BASE DE DONNÉES (MySQL)

Schéma relationnel optimisé autour de 4 tables maîtresses :
1. **`users` :** Gère l'authentification et l'identité (pseudos, emails, mots de passe hachés, rôles).
2. **`categories` :** Définit les "salons" de discussion (les franchises NBA). C'est le pont entre un message et le design du Frontend.
3. **`topics` :** La table centrale (Fil d'actualité). Contient le titre, le contenu, la date et le statut (ouvert/fermé). Clés étrangères vers `users` (auteur) et `categories` (équipe). *(NB: Les sujets remplacent les anciens "posts/messages" pour moderniser le projet en format Feed).*
4. **`topic_reactions` :** Table d'association. Relie un `user_id` et un `topic_id` pour gérer les likes/dislikes et empêcher rigoureusement les doubles votes.

---

## 🎤 5. L'ASTUCE ORAL : Le parcours complet d'une donnée

**Si le jury demande : "Expliquez-moi le fonctionnement exact de la publication d'un message."**

1. **Le Navigateur (JS) :** Sur `team.html`, l'utilisateur clique sur "Publier". `app.js` empêche le rechargement, récupère le texte, s'assure qu'un Token JWT est dans le `localStorage`, et lance un `fetch()` (POST) vers `http://localhost:8080/api/topics` en incluant le JWT dans l'en-tête.
2. **L'Intercepteur :** Sur le port 8080, `auth_middleware.go` valide cryptographiquement le JWT. Si OK, on continue.
3. **Le Contrôleur :** `topic_controller.go` lit le JSON entrant et le passe à la couche métier.
4. **Le Métier :** `topic_service.go` vérifie les règles, s'assure que le sujet n'est pas verrouillé, et passe la requête.
5. **La Base (SQL) :** `topic_repository.go` exécute la requête `INSERT INTO topics` avec des paramètres préparés anti-injection.
6. **La Boucle est bouclée :** Le serveur renvoie un code `201 Created`. `app.js` intercepte ce code, vide la zone de saisie, et relance la fonction `chargerSujets()` pour rafraîchir l'affichage instantanément.