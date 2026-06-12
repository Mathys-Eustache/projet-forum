(function () {
    "use strict";
    const slideTimeout = 10000;
    const carousel = document.getElementById('carousel-players');
    
    if (!carousel) return;
    
    const slides = carousel.querySelectorAll('.slide');
    const prev = carousel.querySelector('.carousel-btn.prev');
    const next = carousel.querySelector('.carousel-btn.next');
    
    let currentSlide = 0;
    let intervalId;

    function showSlide(index) {
        // Nettoyer tous les slides
        slides.forEach(slide => slide.classList.remove('active-fade'));
        
        // Calculer l'index valide
        currentSlide = index < 0 ? slides.length - 1 : index >= slides.length ? 0 : index;
        
        // Ajouter la classe au slide courant
        slides[currentSlide].classList.add('active-fade');
    }

    function autoPlay() {
        intervalId = setInterval(() => showSlide(currentSlide + 1), slideTimeout);
    }

    prev.addEventListener('click', () => {
        clearInterval(intervalId);
        showSlide(currentSlide - 1);
        autoPlay();
    });

    next.addEventListener('click', () => {
        clearInterval(intervalId);
        showSlide(currentSlide + 1);
        autoPlay();
    });

    carousel.addEventListener('mouseover', () => clearInterval(intervalId));
    carousel.addEventListener('mouseout', autoPlay);

    let startX = 0;
    carousel.addEventListener('touchstart', (e) => startX = e.touches[0].clientX);
    carousel.addEventListener('touchend', (e) => {
        const endX = e.changedTouches[0].clientX;
        if (startX - endX > 50) showSlide(currentSlide + 1);
        else if (endX - startX > 50) showSlide(currentSlide - 1);
    });

    // Initialiser le premier slide
    showSlide(0);
    autoPlay();
})();

function gererNavbar() {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    const liInscription = document.getElementById('nav-inscription');
    const liConnexion = document.getElementById('nav-connexion');

    if (token && username && username !== "undefined") {
        if (liInscription) liInscription.innerHTML = `<span class="nav-username">${username}</span>`;
        if (liConnexion) liConnexion.innerHTML = `<a href="#" onclick="deconnexion()">Déconnexion</a>`;
    }
}

function deconnexion() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    window.location.reload();
}

function verifierConnexion() {
    const token = localStorage.getItem('token');
    const zoneSaisie = document.getElementById('zone-saisie');
    const msgConnexion = document.getElementById('message-connexion');

    if (!zoneSaisie || !msgConnexion) return;

    if (token) {
        zoneSaisie.style.display = 'flex';
        msgConnexion.style.display = 'none';
    } else {
        zoneSaisie.style.display = 'none';
        msgConnexion.style.display = 'block';
    }
}

function chargerMessages() {
    fetch('http://localhost:8080/api/messages')
        .then(res => res.json())
        .then(messages => {
            const liste = document.getElementById('liste-messages');
            if (!liste) return;
            
            liste.innerHTML = '';
            if (!Array.isArray(messages)) return;

            const currentUsername = localStorage.getItem('username');

            messages.forEach(msg => {
                // J'ai modifié la structure HTML ici pour intégrer le bouton dans l'en-tête
                let html = `
                    <li>
                        <div class="message-header" style="display: flex; justify-content: space-between; align-items: center;">
                            <div>
                                <span class="message-author">${msg.pseudo || 'Anonyme'}</span>
                                <span class="message-date">${msg.heure || ''}</span>
                            </div>
                `;

                // Le bouton s'affiche désormais à droite du pseudo/date
                if (currentUsername && msg.pseudo === currentUsername) {
                    html += `<button onclick="supprimerMessage(${msg.id})" class="btn-delete" style="background: transparent; color: #dc3545; font-size: 1.1rem; padding: 0; margin-left: 15px; cursor: pointer; border: none;" title="Supprimer ce message">🗑️</button>`;
                }

                html += `
                        </div>
                        <div class="message-content">${msg.texte || ''}</div>
                    </li>
                `;
                liste.innerHTML += html;
            });
            liste.scrollTop = liste.scrollHeight;
        })
        .catch(err => console.error(err));
}

function envoyerMessage() {
    const texte = document.getElementById('nouveau-message').value;
    const token = localStorage.getItem('token');
    let pseudo = localStorage.getItem('username');

    if (!texte || !token) return;

    if (!pseudo || pseudo === "undefined") {
        pseudo = "Anonyme";
    }

    fetch('http://localhost:8080/api/messages', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({ texte: texte, pseudo: pseudo })
    })
    .then(() => {
        document.getElementById('nouveau-message').value = '';
        chargerMessages();
    })
    .catch(err => console.error(err));
}

function supprimerMessage(id) {
    const token = localStorage.getItem('token');
    const pseudo = localStorage.getItem('username'); // Récupération du pseudo manquante
    
    // Si on n'a pas de token ou de pseudo, on bloque
    if (!token || !pseudo) {
        alert("Vous devez être connecté pour supprimer un message.");
        return;
    }

    // Demande de confirmation avant de supprimer (c'est plus sûr)
    if (!confirm("Voulez-vous vraiment supprimer ce message ?")) {
        return;
    }

    fetch(`http://localhost:8080/api/messages/${id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + token,
            'Pseudo': pseudo // Envoi du pseudo dans l'en-tête pour la sécurité côté serveur
        }
    })
    .then(res => {
        if (res.ok) {
            chargerMessages(); // Si c'est ok, on recharge
        } else if (res.status === 403) {
            alert("Vous n'êtes pas autorisé à supprimer ce message.");
        } else {
            alert("Une erreur est survenue lors de la suppression.");
        }
    })
    .catch(err => console.error(err));
}

document.addEventListener('DOMContentLoaded', () => {
    gererNavbar();
    verifierConnexion();
    if (document.getElementById('liste-messages')) {
        chargerMessages();
    }
});