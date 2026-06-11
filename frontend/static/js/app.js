(function () {
    "use strict";
    const slideTimeout = 10000;  // 10 secondes au lieu de 5

    // Configuration pour chaque carrousel
    const carousels = [
        {
            containerId: 'carousel-threads',
            isSlide: true // Effet de glissement
        },
        {
            containerId: 'carousel-players',
            isSlide: false // Effet de fondu
        }
    ];

    // Fonction générique pour initialiser un carrousel
    function initCarousel(config) {
        const carouselContainer = document.getElementById(config.containerId);
        if (!carouselContainer) {
            console.warn(`Carousel container #${config.containerId} not found.`);
            return;
        }

        const slides = Array.from(carouselContainer.querySelectorAll('.slide'));
        if (slides.length === 0) {
            console.warn(`No slides found in #${config.containerId}`);
            return;
        }

        const prev = carouselContainer.querySelector('.carousel-btn.prev');
        const next = carouselContainer.querySelector('.carousel-btn.next');
        const dotsContainer = carouselContainer.querySelector('.carousel-dots');
        const carouselInner = carouselContainer.querySelector('.carousel-inner');
function showSlide(index) {
  slides.forEach((slide, i) => {
    slide.classList.toggle('active', i === index);
  });
  dots.forEach((dot, i) => {
    dot.classList.toggle('active', i === index);
  });
  currentSlide = index;
}

        let currentSlide = 0;
        let intervalId;

        // Générer les points si le conteneur est vide
        if (!dotsContainer.hasChildNodes()) {
            for (let i = 0; i < slides.length; i++) {
                const dot = document.createElement('button');
                dot.className = `dot ${i === 0 ? 'active' : ''}`;
                dot.setAttribute('data-slide-id', i);
                dot.setAttribute('aria-label', `Slide ${i + 1}`);
                dot.addEventListener('click', () => showSlide(i));
                dotsContainer.appendChild(dot);
            }
        }

        const dots = Array.from(dotsContainer.querySelectorAll('.dot'));

        function showSlide(index) {
            // Navigation circulaire
            if (index >= slides.length) {
                index = 0;
            } else if (index < 0) {
                index = slides.length - 1;
            }

            currentSlide = index;

            if (config.isSlide) {
                // Effet de glissement pour threads
                carouselInner.style.transform = `translateX(-${currentSlide * 100}%)`;
            } else {
                // Effet de fondu pour players
                slides.forEach((slide, i) => {
                    slide.classList.toggle('active-fade', i === currentSlide);
                });
            }

            // Mettre à jour les points
            dots.forEach((dot, i) => {
                dot.classList.toggle('active', i === currentSlide);
            });

            // Réinitialiser le minuteur automatique
            clearInterval(intervalId);
            intervalId = setInterval(() => showSlide(currentSlide + 1), slideTimeout);
        }

        // Événements de navigation
        prev.addEventListener('click', () => showSlide(currentSlide - 1));
        next.addEventListener('click', () => showSlide(currentSlide + 1));

        // Pause au survol
        carouselContainer.addEventListener('mouseenter', () => {
            clearInterval(intervalId);
        });

        carouselContainer.addEventListener('mouseleave', () => {
            intervalId = setInterval(() => showSlide(currentSlide + 1), slideTimeout);
        });

        // Gestion du swipe tactile
        let startX = 0;
        carouselContainer.addEventListener('touchstart', (e) => {
            startX = e.touches[0].clientX;
        });

        carouselContainer.addEventListener('touchend', (e) => {
            const endX = e.changedTouches[0].clientX;
            const diff = startX - endX;

            if (diff > 50) {
                showSlide(currentSlide + 1);
            } else if (diff < -50) {
                showSlide(currentSlide - 1);
            }
        });

        // Affichage initial
        showSlide(0);
    }

    // Initialiser tous les carrousels
    carousels.forEach(initCarousel);

})();
if (slides.length > 0) {
  showSlide(0);
  if (nextBtn) nextBtn.addEventListener('click', nextSlide);
  if (prevBtn) prevBtn.addEventListener('click', prevSlide);
  dots.forEach((dot) => {
    dot.addEventListener('click', () => {
      showSlide(Number(dot.dataset.slide));
    });
  });
  setInterval(nextSlide, 6700);
}

function gererNavbar() {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    const liInscription = document.getElementById('nav-inscription');
    const liConnexion = document.getElementById('nav-connexion');

    if (token && username && username !== "undefined") {
        if (liInscription) {
            liInscription.innerHTML = `<span class="nav-username">${username}</span>`;
        }
        if (liConnexion) {
            liConnexion.innerHTML = `<a href="#" onclick="deconnexion()">Déconnexion</a>`;
        }
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
