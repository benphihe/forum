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

    document.querySelector('#fenetreInscription form').addEventListener('submit', async function(e) {
        e.preventDefault();
        const fullname = document.getElementById('fullName').value;
        const email = document.getElementById('emailInscription').value;
        const password = document.getElementById('passwordInscription').value;

        const response = await fetch('/inscription', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ fullname, email, password })
        });

        if (response.ok) {
            alert('Inscription réussie!');
            document.getElementById('fenetreInscription').style.display = 'none';
            document.getElementById('popupOverlay').style.display = 'none';
        } else {
            const errorData = await response.json();
            alert(`Erreur lors de l'inscription: ${errorData.error}`);
        }
    });

    document.querySelector('#fenetreConnexion form').addEventListener('submit', async function(e) {
        e.preventDefault();
        const email = document.getElementById('emailConnexion').value;
        const password = document.getElementById('passwordConnexion').value;

        const response = await fetch('/connexion', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();
            alert('Connexion réussie!');
            document.getElementById('fenetreConnexion').style.display = 'none';
            document.getElementById('popupOverlay').style.display = 'none';
        } else {
            const errorData = await response.json();
            alert(`Erreur lors de la connexion: ${errorData.message}`);
        }
    });
});









