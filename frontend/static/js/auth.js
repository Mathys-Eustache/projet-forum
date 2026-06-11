const formInscription = document.getElementById('form-inscription');
if (formInscription) {
    formInscription.addEventListener('submit', (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const errDiv = document.getElementById('msg-erreur');
        const succDiv = document.getElementById('msg-succes');

        errDiv.style.display = 'none';
        succDiv.style.display = 'none';

        fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        })
        .then(res => res.json().then(data => ({ status: res.status, body: data })))
        .then(res => {
            if (res.status === 201) {
                succDiv.innerText = "Inscription réussie ! Redirection...";
                succDiv.style.display = 'block';
                setTimeout(() => { window.location.href = '/connexion'; }, 2000);
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

        fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        })
        .then(res => res.json().then(data => ({ status: res.status, body: data })))
        .then(res => {
            if (res.status === 200) {
                localStorage.setItem('token', res.body.token);
                localStorage.setItem('username', res.body.username); // <-- LA LIGNE MAGIQUE EST ICI
                
                succDiv.innerText = "Connexion réussie ! Redirection...";
                succDiv.style.display = 'block';
                setTimeout(() => { window.location.href = '/'; }, 2000);
            } else {
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