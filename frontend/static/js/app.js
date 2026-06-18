/**
 * NBA TalkZone - Script Principal
 * Ce fichier gère l'interface utilisateur, la communication avec l'API Go,
 * et la logique d'affichage dynamique (Carrousel, Thèmes, Sujets).
 */

(function () {
    "use strict";
    
    // --- 1. GESTION DU CARROUSEL ---
    const slideTimeout = 10000;
    const carousels = [
        { containerId: 'carousel-threads', isSlide: true },
        { containerId: 'carousel-players', isSlide: false }
    ];

    function initCarousel(config) {
        const carouselContainer = document.getElementById(config.containerId);
        if (!carouselContainer) return;
        
        const slides = Array.from(carouselContainer.querySelectorAll('.slide'));
        if (!slides.length) return;
        
        const prev = carouselContainer.querySelector('.carousel-btn.prev');
        const next = carouselContainer.querySelector('.carousel-btn.next');
        const dotsContainer = carouselContainer.querySelector('.carousel-dots');
        const carouselInner = carouselContainer.querySelector('.carousel-inner');
        
        let currentSlide = 0;
        let intervalId;

        if (dotsContainer && !dotsContainer.hasChildNodes()) {
            slides.forEach((_, i) => {
                const dot = document.createElement('button');
                dot.className = `dot ${i === 0 ? 'active' : ''}`;
                dot.onclick = () => showSlide(i);
                dotsContainer.appendChild(dot);
            });
        }
        
        const dots = dotsContainer ? Array.from(dotsContainer.querySelectorAll('.dot')) : [];
        
        function showSlide(index) {
            currentSlide = (index >= slides.length) ? 0 : (index < 0 ? slides.length - 1 : index);
            
            if (config.isSlide) {
                if (carouselInner) carouselInner.style.transform = `translateX(-${currentSlide * 100}%)`;
            } else {
                slides.forEach((slide, i) => slide.classList.toggle('active-fade', i === currentSlide));
            }
            
            if (dots.length > 0) {
                dots.forEach((dot, i) => dot.classList.toggle('active', i === currentSlide));
            }
            
            clearInterval(intervalId);
            intervalId = setInterval(() => showSlide(currentSlide + 1), slideTimeout);
        }
        
        if (prev) prev.onclick = () => showSlide(currentSlide - 1);
        if (next) next.onclick = () => showSlide(currentSlide + 1);
        carouselContainer.onmouseenter = () => clearInterval(intervalId);
        carouselContainer.onmouseleave = () => intervalId = setInterval(() => showSlide(currentSlide + 1), slideTimeout);
        
        showSlide(0); 
    }
    carousels.forEach(initCarousel);
})();


// --- 2. GESTION DES THÈMES ET DES FRANCHISES ---

function obtenirIdCategorie() {
    return parseInt(new URLSearchParams(window.location.search).get('id')) || 1;
}

const franchises = {
    1: { nom: "Boston Celtics", theme: "theme-celtics" },
    2: { nom: "Oklahoma City Thunder", theme: "theme-thunder" },
    3: { nom: "San Antonio Spurs", theme: "theme-spurs" },
    4: { nom: "Denver Nuggets", theme: "theme-nuggets" },
    5: { nom: "Los Angeles Lakers", theme: "theme-lakers" },
    6: { nom: "Houston Rockets", theme: "theme-rockets" },
    7: { nom: "Minnesota Timberwolves", theme: "theme-timberwolves" },
    8: { nom: "Portland Trail Blazers", theme: "theme-blazers" },
    9: { nom: "Phoenix Suns", theme: "theme-suns" },
    10: { nom: "LA Clippers", theme: "theme-clippers" },
    11: { nom: "Golden State Warriors", theme: "theme-warriors" },
    12: { nom: "New Orleans Pelicans", theme: "theme-pelicans" },
    13: { nom: "Dallas Mavericks", theme: "theme-mavericks" },
    14: { nom: "Memphis Grizzlies", theme: "theme-grizzlies" },
    15: { nom: "Sacramento Kings", theme: "theme-kings" },
    16: { nom: "Utah Jazz", theme: "theme-jazz" },
    17: { nom: "Detroit Pistons", theme: "theme-pistons" },
    18: { nom: "New York Knicks", theme: "theme-knicks" },
    19: { nom: "Cleveland Cavaliers", theme: "theme-cavaliers" },
    20: { nom: "Toronto Raptors", theme: "theme-raptors" },
    21: { nom: "Atlanta Hawks", theme: "theme-hawks" },
    22: { nom: "Philadelphia 76ers", theme: "theme-76ers" },
    23: { nom: "Orlando Magic", theme: "theme-magic" },
    24: { nom: "Charlotte Hornets", theme: "theme-hornets" },
    25: { nom: "Miami Heat", theme: "theme-heat" },
    26: { nom: "Milwaukee Bucks", theme: "theme-bucks" },
    27: { nom: "Chicago Bulls", theme: "theme-bulls" },
    28: { nom: "Brooklyn Nets", theme: "theme-nets" },
    29: { nom: "Indiana Pacers", theme: "theme-pacers" },
    30: { nom: "Washington Wizards", theme: "theme-wizards" }
};

function appliquerThemeDynamique() {
    const id = obtenirIdCategorie();
    const franchise = franchises[id];
    
    if (franchise) {
        document.body.className = franchise.theme;
        const titleElement = document.getElementById('team-title');
        if (titleElement) {
            titleElement.innerText = `Forum - ${franchise.nom}`;
        }
    }
}


// --- 3. GESTION DE LA SESSION ---

function gererNavbar() {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    
    if (token && username && username !== "undefined") {
        const navInsc = document.getElementById('nav-inscription');
        const navConn = document.getElementById('nav-connexion');
        if (navInsc) navInsc.innerHTML = `<span class="nav-username">${username}</span>`;
        if (navConn) navConn.innerHTML = `<a href="#" onclick="deconnexion()">Déconnexion</a>`;
    }
}

function deconnexion() {
    localStorage.clear();
    window.location.reload();
}

function verifierConnexion() {
    const token = localStorage.getItem('token');
    const zoneSaisie = document.getElementById('zone-saisie');
    const msgConnexion = document.getElementById('message-connexion');

    if (zoneSaisie && msgConnexion) {
        if (token && token !== "null" && token !== "undefined") {
            zoneSaisie.classList.remove('hidden-element');
            zoneSaisie.classList.add('flex-element');
            msgConnexion.classList.add('hidden-element');
        } else {
            zoneSaisie.classList.remove('flex-element');
            zoneSaisie.classList.add('hidden-element');
            msgConnexion.classList.remove('hidden-element');
        }
    }
}

function getAuthHeaders() {
    return {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
        'Pseudo': localStorage.getItem('username')
    };
}

function handleFetchResponse(res, successCallback) {
    if (res.status === 401) return deconnexion();
    if (res.ok || res.status === 201) return successCallback();
    res.text().then(text => alert("Erreur : " + text));
}


// --- 4. GESTION DES SUJETS (CRUD) ---

function creerSujet() {
    const title = document.getElementById('topic-title').value;
    const content = document.getElementById('topic-content').value;

    if (!title || !content) return alert("Attention : Le titre et le contenu sont obligatoires !");
    if (!localStorage.getItem('token')) return alert("Erreur : Vous devez être connecté.");

    fetch('http://localhost:8080/api/topics', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ title, content, category_id: obtenirIdCategorie() })
    })
    .then(res => handleFetchResponse(res, () => window.location.reload()))
    .catch(console.error);
}

let currentTopicPage = 1;
const topicsPerPage = 10;

function chargerSujets(page = 1) {
    const targetPage = parseInt(page) || 1;
    if (targetPage < 1) return;
    currentTopicPage = targetPage;
    
    const offset = (targetPage - 1) * topicsPerPage;
    const conteneur = document.getElementById('liste-sujets');
    const pagination = document.getElementById('pagination-sujets');
    if (!conteneur) return;

    const searchInput = document.getElementById('search-topic');
    const searchValue = searchInput ? searchInput.value : '';
    const sortInput = document.getElementById('sort-topic');
    const sortValue = sortInput ? sortInput.value : 'newest';

    fetch(`http://localhost:8080/api/topics?category=${obtenirIdCategorie()}&limit=${topicsPerPage}&offset=${offset}&search=${encodeURIComponent(searchValue)}&sort=${sortValue}`)
    .then(res => res.json())
    .then(topics => {
        
        const compteur = document.getElementById('total-messages');
        if (compteur) {
            compteur.innerText = Array.isArray(topics) ? topics.length : 0;
        }

        if (!Array.isArray(topics) || !topics.length) {
            if (targetPage === 1) {
                conteneur.innerHTML = '<p class="empty-topics-msg">Aucun sujet trouvé.</p>';
                if (pagination) pagination.innerHTML = '';
            } else {
                if (pagination) {
                    pagination.innerHTML = `
                        <div class="pagination-container">
                            <button class="pagination-btn" onclick="chargerSujets(${targetPage - 1})">Précédent</button>
                            <span class="pagination-current">Page ${targetPage}</span>
                            <button class="pagination-btn btn-disabled" disabled>Suivant</button>
                        </div>
                    `;
                }
            }
            return;
        }

        const currentUsername = localStorage.getItem('username');
        const isAdmin = false;

        conteneur.innerHTML = topics.map(t => {
            const statusClass = t.status === 'fermé' ? 'status-ferme' : 'status-ouvert';
            const statusText = t.status === 'fermé' ? '🔒 Fermé' : '🔓 Ouvert';
            
            return `
            <div class="sujet-card">
                <div class="sujet-actions">
                    <span class="sujet-status ${statusClass}">
                        ${statusText}
                    </span>
                    ${(currentUsername === t.author || isAdmin) ? `
                    <button class="action-btn" onclick="toggleStatusSujet(${t.id}, '${t.status}')" title="Changer l'état">⚙️</button>
                    <button class="action-btn" onclick="editerSujet(${t.id}, this)" title="Modifier">✏️</button>
                    <button class="action-btn action-delete" onclick="supprimerSujet(${t.id})" title="Supprimer">🗑️</button>
                    ` : ''}
                </div>
                <h3 class="sujet-title">${t.title}</h3>
                <p class="topic-content">${t.content}</p>
                <div class="sujet-footer">
                    <small class="sujet-meta">Par <strong>${t.author}</strong> le ${t.created_at}</small>
                    <div class="sujet-reactions">
                        <button class="reaction-btn" onclick="reagirSujet(${t.id}, 'like')">👍 ${t.likes || 0}</button>
                        <button class="reaction-btn" onclick="reagirSujet(${t.id}, 'dislike')">👎 ${t.dislikes || 0}</button>
                    </div>
                </div>
            </div>
            `;
        }).join('');

        if (pagination) {
            pagination.innerHTML = `
                <div class="pagination-container">
                    <button class="pagination-btn ${targetPage === 1 ? 'btn-disabled' : ''}" onclick="chargerSujets(${targetPage - 1})" ${targetPage === 1 ? 'disabled' : ''}>Précédent</button>
                    <span class="pagination-current">Page ${targetPage}</span>
                    <button class="pagination-btn ${topics.length < topicsPerPage ? 'btn-disabled' : ''}" onclick="chargerSujets(${targetPage + 1})" ${topics.length < topicsPerPage ? 'disabled' : ''}>Suivant</button>
                </div>
            `;
        }
    })
    .catch(console.error);
}

function toggleStatusSujet(id, currentStatus) {
    const newStatus = currentStatus === 'fermé' ? 'ouvert' : 'fermé';
    fetch(`http://localhost:8080/api/topics/status/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({ status: newStatus })
    })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

function editerSujet(id, btnElement) {
    const pElement = btnElement.closest('.sujet-card').querySelector('.topic-content');
    const oldContent = pElement.innerText;
    const newContent = prompt("Modifier le sujet :", oldContent);
    
    if (newContent === null || newContent.trim() === "" || newContent === oldContent) return;

    fetch(`http://localhost:8080/api/topics/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({ content: newContent })
    })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

function reagirSujet(id, action) {
    if (!localStorage.getItem('token')) return alert("Erreur : Vous devez être connecté pour réagir.");
    
    fetch(`http://localhost:8080/api/topics/react/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({ action: action })
    })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

function supprimerSujet(id) {
    if (!confirm("Voulez-vous vraiment supprimer ce sujet ?")) return;
    
    fetch(`http://localhost:8080/api/topics/${id}`, { method: 'DELETE', headers: getAuthHeaders() })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

// --- FONCTION POUR LE COMPTEUR GLOBAL (Barre Bleue) ---
function chargerVraiCompteur() {
    const compteur = document.getElementById('compteur-reel');
    if (!compteur) return; // Si la barre bleue n'est pas sur cette page, on annule

    // On récupère tous les messages du serveur (sans limite) pour compter le total
    fetch('http://localhost:8080/api/topics?limit=1000')
    .then(res => res.json())
    .then(topics => {
        compteur.innerText = Array.isArray(topics) ? topics.length : 0;
    })
    .catch(console.error);
}

// --- 5. INITIALISATION GLOBALE ---
document.addEventListener('DOMContentLoaded', () => {
    gererNavbar();
    verifierConnexion();
    appliquerThemeDynamique();
    chargerVraiCompteur(); // Lancement du compteur de la barre bleue
    
    if (document.getElementById('liste-sujets')) {
        chargerSujets();
    }
});