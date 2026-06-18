// =========================================================
// GESTION DE L'INSCRIPTION
// =========================================================
const formInscription = document.getElementById('form-inscription');
if (formInscription) {
    // On écoute la soumission du formulaire
    formInscription.addEventListener('submit', (e) => {
        e.preventDefault(); // Empêche le rechargement automatique de la page

        // 1. Récupération des données tapées par l'utilisateur
        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        
        const errDiv = document.getElementById('msg-erreur');
        const succDiv = document.getElementById('msg-succes');

        // Réinitialisation des messages d'alerte
        errDiv.style.display = 'none';
        succDiv.style.display = 'none';

        // 2. Envoi des données à l'API Go au format JSON
        fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        })
        .then(res => res.json().then(data => ({ status: res.status, body: data })))
        .then(res => {
            // 3. Traitement de la réponse du serveur
            if (res.status === 201) { // 201 = Created (Créé avec succès)
                succDiv.innerText = "Inscription réussie ! Redirection...";
                succDiv.style.display = 'block';
                // Redirection vers la page de connexion après 2 secondes
                setTimeout(() => { window.location.href = 'connexion.html'; }, 2000);
            } else {
                errDiv.innerText = res.body.erreur || "Une erreur est survenue";
                errDiv.style.display = 'block';
            }
        })
        .catch(() => {
            errDiv.innerText = "Impossible de joindre le serveur backend";
            errDiv.style.display = 'block';
        });
    });
}

// =========================================================
// GESTION DE LA CONNEXION
// =========================================================
const formConnexion = document.getElementById('form-connexion');
if (formConnexion) {
    formConnexion.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const errDiv = document.getElementById('msg-erreur');
        const succDiv = document.getElementById('msg-succes');

        errDiv.style.display = 'none';
        succDiv.style.display = 'none';

        // 1. Appel à l'API pour vérifier les identifiants
        fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        })
        .then(res => res.json().then(data => ({ status: res.status, body: data })))
        .then(res => {
            // 2. Si la connexion est validée par le backend
            if (res.status === 200) {
                // IMPORTANT : Stockage du "Token JWT" et des infos dans le navigateur
                // C'est ce qui permet à l'utilisateur de rester connecté sur toutes les pages
                localStorage.setItem('token', res.body.token);
                localStorage.setItem('username', res.body.username);
                localStorage.setItem('role', res.body.role || 'user'); 
                
                succDiv.innerText = "Connexion réussie ! Redirection...";
                succDiv.style.display = 'block';
                // Redirection vers l'accueil
                setTimeout(() => { window.location.href = 'index.html'; }, 2000);
            } else {
                // Affiche l'erreur renvoyée par le backend (ex: "identifiants invalides")
                errDiv.innerText = res.body.erreur || "Identifiants incorrects";
                errDiv.style.display = 'block';
            }
        })
        .catch(() => {
            errDiv.innerText = "Impossible de joindre le serveur backend";
            errDiv.style.display = 'block';
        });
    });
}