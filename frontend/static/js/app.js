(function () {
    "use strict";
    // Carousels et initialisation (code d'origine conservé)
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

        if (!dotsContainer.hasChildNodes()) {
            slides.forEach((_, i) => {
                const dot = document.createElement('button');
                dot.className = `dot ${i === 0 ? 'active' : ''}`;
                dot.onclick = () => showSlide(i);
                dotsContainer.appendChild(dot);
            });
        }
        const dots = Array.from(dotsContainer.querySelectorAll('.dot'));
        function showSlide(index) {
            currentSlide = (index >= slides.length) ? 0 : (index < 0 ? slides.length - 1 : index);
            if (config.isSlide) {
                carouselInner.style.transform = `translateX(-${currentSlide * 100}%)`;
            } else {
                slides.forEach((slide, i) => slide.classList.toggle('active-fade', i === currentSlide));
            }
            dots.forEach((dot, i) => dot.classList.toggle('active', i === currentSlide));
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

function obtenirIdCategorie() {
    return parseInt(new URLSearchParams(window.location.search).get('id')) || 1;
}

const franchises = {
    1: { nom: "Boston Celtics 🍀", theme: "theme-celtics" },
    2: { nom: "Oklahoma City Thunder ⚡", theme: "theme-thunder" },
    3: { nom: "San Antonio Spurs 🤠", theme: "theme-spurs" },
    4: { nom: "Denver Nuggets ⛏️", theme: "theme-nuggets" },
    5: { nom: "Los Angeles Lakers 🎬", theme: "theme-lakers" },
    6: { nom: "Houston Rockets 🚀", theme: "theme-rockets" },
    7: { nom: "Minnesota Timberwolves 🐺", theme: "theme-timberwolves" },
    8: { nom: "Portland Trail Blazers 🌲", theme: "theme-blazers" },
    9: { nom: "Phoenix Suns ☀️", theme: "theme-suns" },
    10: { nom: "LA Clippers ⛵", theme: "theme-clippers" },
    11: { nom: "Golden State Warriors 🌉", theme: "theme-warriors" },
    12: { nom: "New Orleans Pelicans ⚜️", theme: "theme-pelicans" },
    13: { nom: "Dallas Mavericks 🐴", theme: "theme-mavericks" },
    14: { nom: "Memphis Grizzlies 🐻", theme: "theme-grizzlies" },
    15: { nom: "Sacramento Kings 👑", theme: "theme-kings" },
    16: { nom: "Utah Jazz 🎷", theme: "theme-jazz" },
    17: { nom: "Detroit Pistons ⚙️", theme: "theme-pistons" },
    18: { nom: "New York Knicks 🗽", theme: "theme-knicks" },
    19: { nom: "Cleveland Cavaliers ⚔️", theme: "theme-cavaliers" },
    20: { nom: "Toronto Raptors 🦖", theme: "theme-raptors" },
    21: { nom: "Atlanta Hawks 🦅", theme: "theme-hawks" },
    22: { nom: "Philadelphia 76ers 🔔", theme: "theme-76ers" },
    23: { nom: "Orlando Magic 🪄", theme: "theme-magic" },
    24: { nom: "Charlotte Hornets 🐝", theme: "theme-hornets" },
    25: { nom: "Miami Heat 🔥", theme: "theme-heat" },
    26: { nom: "Milwaukee Bucks 🦌", theme: "theme-bucks" },
    27: { nom: "Chicago Bulls 🐂", theme: "theme-bulls" },
    28: { nom: "Brooklyn Nets 🏙️", theme: "theme-nets" },
    29: { nom: "Indiana Pacers 🏎️", theme: "theme-pacers" },
    30: { nom: "Washington Wizards 🧙‍♂️", theme: "theme-wizards" }
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
        zoneSaisie.style.display = token ? 'flex' : 'none';
        msgConnexion.style.display = token ? 'none' : 'block';
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
        if (!Array.isArray(topics) || !topics.length) {
            if (targetPage === 1) {
                conteneur.innerHTML = '<p style="text-align: center; color: #999;">Aucun message trouvé.</p>';
                if (pagination) pagination.innerHTML = '';
            } else {
                if (pagination) {
                    pagination.innerHTML = `
                        <div style="display: flex; justify-content: center; align-items: center; gap: 15px; margin-top: 20px;">
                            <button onclick="chargerSujets(${targetPage - 1})" style="padding: 8px 16px; border: 1px solid #ccc; background: #fff; border-radius: 4px; cursor: pointer;">Précédent</button>
                            <span style="font-weight: bold; color: #333;">Page ${targetPage}</span>
                            <button disabled style="padding: 8px 16px; border: 1px solid #ccc; background: #eee; color: #999; border-radius: 4px; cursor: not-allowed;">Suivant</button>
                        </div>
                    `;
                }
            }
            return;
        }

        const currentUsername = localStorage.getItem('username');
        conteneur.innerHTML = topics.map(t => `
            <div class="sujet-card" style="border: 1px solid #e0e0e0; padding: 15px; margin-bottom: 15px; border-radius: 5px; background: #fff; position: relative; box-shadow: 0 2px 4px rgba(0,0,0,0.05);">
                <div style="position: absolute; right: 15px; top: 15px;">
                    <span style="font-size: 0.8rem; padding: 3px 8px; border-radius: 4px; margin-right: 5px; font-weight: bold; background: ${t.status === 'fermé' ? '#dc3545' : '#28a745'}; color: white;">
                        ${t.status === 'fermé' ? '🔒 Fermé' : '🔓 Ouvert'}
                    </span>
                    ${currentUsername === t.author ? `
                    <button onclick="toggleStatusSujet(${t.id}, '${t.status}')" title="Changer l'état" style="background: transparent; color: #6c757d; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">⚙️</button>
                    <button onclick="editerSujet(${t.id}, this)" title="Modifier" style="background: transparent; color: #6c757d; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">✏️</button>
                    <button onclick="supprimerSujet(${t.id})" title="Supprimer" style="background: transparent; color: #dc3545; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">🗑️</button>
                    ` : ''}
                </div>
                <h3 onclick="window.location.href='/topic?id=${t.id}'" style="margin-top: 0; cursor: pointer; padding-right: 140px;">${t.title}</h3>
                <p class="topic-content" style="color: #333;">${t.content}</p>
                <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 15px;">
                    <small style="color: #666;">Par <strong>${t.author}</strong> le ${t.created_at}</small>
                    <div>
                        <button onclick="reagirSujet(${t.id}, 'like')" style="background: transparent; border: 1px solid #ccc; border-radius: 4px; padding: 4px 8px; cursor: pointer; margin-right: 5px;">👍 ${t.likes || 0}</button>
                        <button onclick="reagirSujet(${t.id}, 'dislike')" style="background: transparent; border: 1px solid #ccc; border-radius: 4px; padding: 4px 8px; cursor: pointer;">👎 ${t.dislikes || 0}</button>
                    </div>
                </div>
            </div>
        `).join('');

        if (pagination) {
            pagination.innerHTML = `
                <div style="display: flex; justify-content: center; align-items: center; gap: 15px; margin-top: 20px;">
                    <button onclick="chargerSujets(${targetPage - 1})" ${targetPage === 1 ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #ccc; background: ${targetPage === 1 ? '#eee' : '#fff'}; border-radius: 4px; cursor: ${targetPage === 1 ? 'not-allowed' : 'pointer'};">Précédent</button>
                    <span style="font-weight: bold; color: #333;">Page ${targetPage}</span>
                    <button onclick="chargerSujets(${targetPage + 1})" ${topics.length < topicsPerPage ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #ccc; background: ${topics.length < topicsPerPage ? '#eee' : '#fff'}; border-radius: 4px; cursor: ${topics.length < topicsPerPage ? 'not-allowed' : 'pointer'};">Suivant</button>
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
    const newContent = prompt("Modifier votre message :", oldContent);
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
    if (!localStorage.getItem('token')) return alert("Erreur : Vous devez être connecté.");
    fetch(`http://localhost:8080/api/topics/react/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({ action: action })
    })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

function supprimerSujet(id) {
    if (!confirm("Voulez-vous vraiment supprimer ce message ?")) return;
    fetch(`http://localhost:8080/api/topics/${id}`, { method: 'DELETE', headers: getAuthHeaders() })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

document.addEventListener('DOMContentLoaded', () => {
    gererNavbar();
    verifierConnexion();
    appliquerThemeDynamique();
    if (document.getElementById('liste-sujets')) chargerSujets();
});