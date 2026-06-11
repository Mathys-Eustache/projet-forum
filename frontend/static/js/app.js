(function () {
    "use strict";
    const slideTimeout = 10000;

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