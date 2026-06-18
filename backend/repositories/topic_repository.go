package repositories

import (
	"database/sql"
	"errors"
	"projet-forum/backend/models"
)

type TopicRepository struct {
	DB *sql.DB
}

// Structure spécifique pour envoyer les données au frontend (jointure avec la table Users)
type TopicResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
}

// CreateTopic insère un nouveau sujet dans la table Topics
func (r *TopicRepository) CreateTopic(topic models.Topic) error {
	query := `INSERT INTO Topics (title, content, user_id, category_id, status, likes, dislikes) 
			  VALUES (?, ?, ?, ?, 'ouvert', 0, 0)`
	// On utilise Exec() car on ne s'attend pas à un retour de données, juste une confirmation d'insertion
	_, err := r.DB.Exec(query, topic.Title, topic.Content, topic.UserID, topic.CategoryID)
	return err
}

// GetAllTopics récupère les sujets de toutes les catégories confondues avec pagination et tri
func (r *TopicRepository) GetAllTopics(limit int, offset int, search string, sortBy string) ([]TopicResponse, error) {
	// Choix dynamique du tri
	orderClause := "ORDER BY t.created_at DESC"
	if sortBy == "oldest" {
		orderClause = "ORDER BY t.created_at ASC"
	} else if sortBy == "popular" {
		orderClause = "ORDER BY t.likes DESC, t.created_at DESC"
	}

	var rows *sql.Rows
	var err error

	// Construction de la requête avec ou sans filtre de recherche (LIKE)
	if search != "" {
		searchParam := "%" + search + "%"
		query := `SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status, t.likes, t.dislikes 
				  FROM Topics t JOIN Users u ON t.user_id = u.id 
				  WHERE t.title LIKE ? OR t.content LIKE ? ` + orderClause + ` LIMIT ? OFFSET ?`
		rows, err = r.DB.Query(query, searchParam, searchParam, limit, offset)
	} else {
		query := `SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status, t.likes, t.dislikes 
				  FROM Topics t JOIN Users u ON t.user_id = u.id ` + orderClause + ` LIMIT ? OFFSET ?`
		rows, err = r.DB.Query(query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []TopicResponse
	for rows.Next() {
		var t TopicResponse
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Author, &t.CreatedAt, &t.Status, &t.Likes, &t.Dislikes); err != nil {
			continue // Si une ligne plante, on passe à la suivante sans tout bloquer
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// GetTopicsByCategory récupère les sujets spécifiques à UNE seule franchise (categoryID)
func (r *TopicRepository) GetTopicsByCategory(categoryID int, limit int, offset int, search string, sortBy string) ([]TopicResponse, error) {
	orderClause := "ORDER BY t.created_at DESC"
	if sortBy == "oldest" {
		orderClause = "ORDER BY t.created_at ASC"
	} else if sortBy == "popular" {
		orderClause = "ORDER BY t.likes DESC, t.created_at DESC"
	}

	var rows *sql.Rows
	var err error

	if search != "" {
		searchParam := "%" + search + "%"
		query := `SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status, t.likes, t.dislikes 
				  FROM Topics t JOIN Users u ON t.user_id = u.id 
				  WHERE t.category_id = ? AND (t.title LIKE ? OR t.content LIKE ?) ` + orderClause + ` LIMIT ? OFFSET ?`
		rows, err = r.DB.Query(query, categoryID, searchParam, searchParam, limit, offset)
	} else {
		query := `SELECT t.id, t.title, t.content, u.username, DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), t.status, t.likes, t.dislikes 
				  FROM Topics t JOIN Users u ON t.user_id = u.id 
				  WHERE t.category_id = ? ` + orderClause + ` LIMIT ? OFFSET ?`
		rows, err = r.DB.Query(query, categoryID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []TopicResponse
	for rows.Next() {
		var t TopicResponse
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Author, &t.CreatedAt, &t.Status, &t.Likes, &t.Dislikes); err != nil {
			continue
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// DeleteTopic supprime un sujet UNIQUEMENT si l'utilisateur en est l'auteur (ou l'admin via le Service)
func (r *TopicRepository) DeleteTopic(id int, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	// Contrôle de sécurité final
	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("DELETE FROM Topics WHERE id = ?", id)
	return err
}

// UpdateTopic modifie un sujet (contrôle de l'auteur inclus)
func (r *TopicRepository) UpdateTopic(id int, content string, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("UPDATE Topics SET content = ? WHERE id = ?", content, id)
	return err
}

// UpdateTopicStatus permet d'ouvrir ou de fermer un sujet
func (r *TopicRepository) UpdateTopicStatus(id int, status string, pseudo string) error {
	var author string
	err := r.DB.QueryRow("SELECT u.username FROM Topics t JOIN Users u ON t.user_id = u.id WHERE t.id = ?", id).Scan(&author)
	if err != nil {
		return err
	}

	if author != pseudo {
		return errors.New("forbidden")
	}

	_, err = r.DB.Exec("UPDATE Topics SET status = ? WHERE id = ?", status, id)
	return err
}

// ReactTopic gère la logique complexe des "J'aime / Je n'aime pas" (Toggle)
func (r *TopicRepository) ReactTopic(topicID int, action string, pseudo string) error {
	// 1. Récupération de l'ID utilisateur
	var userID int
	err := r.DB.QueryRow("SELECT id FROM Users WHERE username = ?", pseudo).Scan(&userID)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	// 2. Vérification : l'utilisateur a-t-il déjà réagi à ce sujet ?
	var currentReaction string
	err = r.DB.QueryRow("SELECT reaction_type FROM Topic_Reactions WHERE topic_id = ? AND user_id = ?", topicID, userID).Scan(&currentReaction)

	if err == sql.ErrNoRows {
		// Pas de réaction : on l'ajoute
		_, err = r.DB.Exec("INSERT INTO Topic_Reactions (topic_id, user_id, reaction_type) VALUES (?, ?, ?)", topicID, userID, action)
	} else if err == nil {
		// Réaction existante : on l'enlève (si clic sur le même bouton) ou on la modifie
		if currentReaction == action {
			_, err = r.DB.Exec("DELETE FROM Topic_Reactions WHERE topic_id = ? AND user_id = ?", topicID, userID)
		} else {
			_, err = r.DB.Exec("UPDATE Topic_Reactions SET reaction_type = ? WHERE topic_id = ? AND user_id = ?", action, topicID, userID)
		}
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 3. Mise à jour des compteurs totaux dans la table Topics
	updateQuery := `UPDATE Topics SET 
					likes = (SELECT COUNT(*) FROM Topic_Reactions WHERE topic_id = ? AND reaction_type = 'like'), 
					dislikes = (SELECT COUNT(*) FROM Topic_Reactions WHERE topic_id = ? AND reaction_type = 'dislike') 
					WHERE id = ?`
	_, err = r.DB.Exec(updateQuery, topicID, topicID, topicID)

	return err
}
