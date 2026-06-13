(function () {
    "use strict";
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
    if (token && token.startsWith("TOKEN_SIMULE")) return deconnexion();

    const zoneSaisie = document.getElementById('zone-saisie');
    const msgConnexion = document.getElementById('message-connexion');

    if (zoneSaisie && msgConnexion) {
        zoneSaisie.style.display = token ? 'flex' : 'none';
        msgConnexion.style.display = token ? 'none' : 'block';
    }
}

function obtenirIdCategorie() {
    return parseInt(new URLSearchParams(window.location.search).get('id')) || 1;
}

// Extraction propre des entêtes
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

// Création d'un sujet
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

// --- GESTION DES TOPICS (SUJETS) ---
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

    fetch(`http://localhost:8080/api/topics?category=${obtenirIdCategorie()}&limit=${topicsPerPage}&offset=${offset}&search=${encodeURIComponent(searchValue)}`)
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
                            <button onclick="chargerSujets(${targetPage - 1})" style="padding: 8px 16px; border: 1px solid #1d428a; background: #fff; color: #1d428a; border-radius: 4px; cursor: pointer;">Précédent</button>
                            <span style="font-weight: bold; color: #333;">Page ${targetPage}</span>
                            <button disabled style="padding: 8px 16px; border: 1px solid #1d428a; background: #eee; color: #999; border-radius: 4px; cursor: not-allowed;">Suivant</button>
                        </div>
                    `;
                }
            }
            return;
        }

        const currentUsername = localStorage.getItem('username');
        conteneur.innerHTML = topics.map(t => `
            <div class="sujet-card" style="border: 1px solid #ccc; padding: 15px; margin-bottom: 15px; border-radius: 5px; background: #fff; position: relative;">
                <div style="position: absolute; right: 15px; top: 15px;">
                    <span style="font-size: 0.8rem; padding: 3px 8px; border-radius: 4px; margin-right: 5px; font-weight: bold; background: ${t.status === 'fermé' ? '#dc3545' : '#28a745'}; color: white;">
                        ${t.status === 'fermé' ? '🔒 Fermé' : '🔓 Ouvert'}
                    </span>
                    ${currentUsername === t.author ? `
                    <button onclick="toggleStatusSujet(${t.id}, '${t.status}')" title="Changer l'état (Ouvert/Fermé)" style="background: transparent; color: #6c757d; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">⚙️</button>
                    <button onclick="editerSujet(${t.id}, this)" title="Modifier" style="background: transparent; color: #1d428a; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">✏️</button>
                    <button onclick="supprimerSujet(${t.id})" title="Supprimer" style="background: transparent; color: #dc3545; font-size: 1.2rem; padding: 0 5px; cursor: pointer; border: none;">🗑️</button>
                    ` : ''}
                </div>
                <h3 onclick="window.location.href='/topic.html?id=${t.id}'" style="margin-top: 0; color: #1d428a; cursor: pointer; padding-right: 140px;">${t.title}</h3>
                <p class="topic-content" style="color: #333;">${t.content}</p>
                <small style="color: #666;">Par <strong>${t.author}</strong> le ${t.created_at}</small>
            </div>
        `).join('');

        if (pagination) {
            pagination.innerHTML = `
                <div style="display: flex; justify-content: center; align-items: center; gap: 15px; margin-top: 20px;">
                    <button onclick="chargerSujets(${targetPage - 1})" ${targetPage === 1 ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #1d428a; background: ${targetPage === 1 ? '#eee' : '#fff'}; color: ${targetPage === 1 ? '#999' : '#1d428a'}; border-radius: 4px; cursor: ${targetPage === 1 ? 'not-allowed' : 'pointer'};">Précédent</button>
                    <span style="font-weight: bold; color: #333;">Page ${targetPage}</span>
                    <button onclick="chargerSujets(${targetPage + 1})" ${topics.length < topicsPerPage ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #1d428a; background: ${topics.length < topicsPerPage ? '#eee' : '#fff'}; color: ${topics.length < topicsPerPage ? '#999' : '#1d428a'}; border-radius: 4px; cursor: ${topics.length < topicsPerPage ? 'not-allowed' : 'pointer'};">Suivant</button>
                </div>
            `;
        }
    })
    .catch(console.error);
}

// Fonction pour basculer entre Ouvert et Fermé
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

function supprimerSujet(id) {
    if (!confirm("Voulez-vous vraiment supprimer ce message ?")) return;
    fetch(`http://localhost:8080/api/topics/${id}`, { method: 'DELETE', headers: getAuthHeaders() })
    .then(res => handleFetchResponse(res, () => chargerSujets(currentTopicPage)))
    .catch(console.error);
}

// --- GESTION DE LA PAGE SUJET DÉTAILLÉE ---
function chargerPageSujet() {
    const topicId = new URLSearchParams(window.location.search).get('id');
    if (!topicId) return;

    fetch(`http://localhost:8080/api/topics`)
    .then(res => res.json())
    .then(topics => {
        const topic = topics.find(t => t.id == topicId);
        document.getElementById('topic-page-title').innerText = topic ? topic.title : "Sujet introuvable";
        if (topic) {
            document.getElementById('topic-page-content').innerText = topic.content;
            document.getElementById('topic-page-meta').innerHTML = `Par <strong>${topic.author}</strong> le ${topic.created_at}`;
            
            // --- BLOCAGE FT-3 SI LE SUJET EST FERMÉ ---
            if (topic.status === 'fermé') {
                const textarea = document.getElementById('nouvelle-reponse');
                const btnEnvoyer = document.querySelector('button[onclick="envoyerReponse()"]');
                if (textarea) {
                    textarea.disabled = true;
                    textarea.placeholder = "🔒 Ce sujet est fermé. Impossible de répondre.";
                }
                if (btnEnvoyer) {
                    btnEnvoyer.style.display = 'none';
                }
            }
        }
    })
    .catch(console.error);

    chargerReponses(topicId);
}

let currentPostPage = 1;
const postsPerPage = 10;

function chargerReponses(topicId, page = 1) {
    const targetPostPage = parseInt(page) || 1;
    if (targetPostPage < 1) return;
    currentPostPage = targetPostPage;
    
    const offset = (targetPostPage - 1) * postsPerPage;

    fetch(`http://localhost:8080/api/posts?topic_id=${topicId}&limit=${postsPerPage}&offset=${offset}`)
    .then(res => res.json())
    .then(posts => {
        const conteneur = document.getElementById('liste-reponses');
        const pagination = document.getElementById('pagination-reponses');
        if (!conteneur) return;
        
        if (!Array.isArray(posts) || !posts.length) {
            if (targetPostPage === 1) conteneur.innerHTML = '<p style="text-align: center; color: #999;">Aucune réponse pour le moment. Sois le premier !</p>';
            if (pagination && targetPostPage === 1) pagination.innerHTML = '';
            return;
        }
        
        const currentUsername = localStorage.getItem('username');
        conteneur.innerHTML = posts.map(p => `
            <div style="border-bottom: 1px solid #eee; padding: 15px 0; overflow: hidden;">
                ${currentUsername === p.author ? `<button onclick="supprimerReponse(${p.id}, ${topicId})" style="background: transparent; color: #dc3545; font-size: 1.1rem; padding: 0; cursor: pointer; border: none; float: right;">🗑️</button>` : ''}
                <p style="margin: 0 0 10px 0; color: #333;">${p.content}</p>
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <small style="color: #666;">Par <strong>${p.author}</strong> le ${p.created_at}</small>
                </div>
            </div>
        `).join('');

        if (pagination) {
            pagination.innerHTML = `
                <div style="display: flex; justify-content: center; align-items: center; gap: 15px; margin-top: 20px;">
                    <button onclick="chargerReponses(${topicId}, ${targetPostPage - 1})" ${targetPostPage === 1 ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #1d428a; background: ${targetPostPage === 1 ? '#eee' : '#fff'}; color: ${targetPostPage === 1 ? '#999' : '#1d428a'}; border-radius: 4px; cursor: ${targetPostPage === 1 ? 'not-allowed' : 'pointer'};">Précédent</button>
                    <span style="font-weight: bold; color: #333;">Page ${targetPostPage}</span>
                    <button onclick="chargerReponses(${topicId}, ${targetPostPage + 1})" ${posts.length < postsPerPage ? 'disabled' : ''} style="padding: 8px 16px; border: 1px solid #1d428a; background: ${posts.length < postsPerPage ? '#eee' : '#fff'}; color: ${posts.length < postsPerPage ? '#999' : '#1d428a'}; border-radius: 4px; cursor: ${posts.length < postsPerPage ? 'not-allowed' : 'pointer'};">Suivant</button>
                </div>
            `;
        }
    })
    .catch(console.error);
}

function supprimerReponse(id, topicId) {
    if (!confirm("Voulez-vous vraiment supprimer cette réponse ?")) return;
    fetch(`http://localhost:8080/api/posts/${id}`, { method: 'DELETE', headers: getAuthHeaders() })
    .then(res => handleFetchResponse(res, () => chargerReponses(topicId, currentPostPage)))
    .catch(console.error);
}

function envoyerReponse() {
    const content = document.getElementById('nouvelle-reponse').value;
    const topicId = new URLSearchParams(window.location.search).get('id');

    if (!content || !topicId) return;
    if (!localStorage.getItem('token')) return alert("Erreur : Vous devez être connecté.");

    fetch(`http://localhost:8080/api/posts?topic_id=${topicId}`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ content })
    })
    .then(res => handleFetchResponse(res, () => {
        document.getElementById('nouvelle-reponse').value = '';
        chargerReponses(topicId, 1);
    }))
    .catch(console.error);
}

document.addEventListener('DOMContentLoaded', () => {
    gererNavbar();
    verifierConnexion();
    if (document.getElementById('liste-sujets')) chargerSujets();
    if (document.getElementById('sujet-principal')) chargerPageSujet();
});