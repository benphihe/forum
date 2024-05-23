document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('openInscriptionPopup').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreInscription').style.display = 'block';
        document.getElementById('popupOverlay').style.display = 'block';
    });

    document.getElementById('closeInscriptionPopup').addEventListener('click', function() {
        document.getElementById('fenetreInscription').style.display = 'none';
        document.getElementById('popupOverlay').style.display = 'none';
    });

    document.getElementById('openConnexionPopupFromInscription').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreInscription').style.display = 'none';
        document.getElementById('fenetreConnexion').style.display = 'block';
        document.getElementById('popupOverlay').style.display = 'block';
    });

    document.getElementById('closeConnexionPopup').addEventListener('click', function() {
        document.getElementById('fenetreConnexion').style.display = 'none';
        document.getElementById('popupOverlay').style.display = 'none';
    });

    document.getElementById('openInscriptionPopupFromConnexion').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreConnexion').style.display = 'none';
        document.getElementById('fenetreInscription').style.display = 'block';
        document.getElementById('popupOverlay').style.display = 'block';
    });

    document.querySelector('#formInscription').addEventListener('submit', async function(e) {
        e.preventDefault();
        const fullName = document.getElementById('fullName').value;
        const email = document.getElementById('emailInscriptionPopup').value;
        const password = document.getElementById('passwordInscriptionPopup').value;

        const response = await fetch('http://localhost:5000/inscription', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ fullname: fullName, email, password })
        });

        if (response.ok) {
            alert('Inscription réussie!');
            document.getElementById('fenetreInscription').style.display = 'none';
            document.getElementById('popupOverlay').style.display = 'none';
        } else {
            alert('Erreur lors de l\'inscription');
        }
    });

    document.querySelector('#formConnexion').addEventListener('submit', async function(e) {
        e.preventDefault();
        const email = document.getElementById('emailConnexion').value;
        const password = document.getElementById('passwordConnexion').value;

        const response = await fetch('http://localhost:5000/connexion', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();
            alert('Connexion réussie!');
            localStorage.setItem('user', JSON.stringify(data.user));
            document.getElementById('fenetreConnexion').style.display = 'none';
            document.getElementById('popupOverlay').style.display = 'none';
        } else {
            alert('Erreur lors de la connexion');
        }
    });
});






