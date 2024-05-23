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
            try {
                const data = await response.json();
                alert('Inscription réussie!');
                document.getElementById('fenetreInscription').style.display = 'none';
                document.getElementById('popupOverlay').style.display = 'none';
            } catch (error) {
                console.error('Erreur lors de la conversion de la réponse en JSON:', error);
            }
        } else {
            console.error('Erreur lors de la requête:', response.status, response.statusText);
        }
    });

    document.querySelector('#fenetreConnexion form').addEventListener('submit', async function(e) {
        e.preventDefault();
        const email = document.getElementById('emailConnexion').value;
        const password = document.getElementById('passwordConnexion').value;

        try {
            const response = await fetch('/connexion', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            if (response.ok) {
                if (response.headers.get('content-type').includes('application/json')) {
                    const data = await response.json();
                    alert('Connexion réussie!');
                    document.getElementById('fenetreConnexion').style.display = 'none';
                    document.getElementById('popupOverlay').style.display = 'none';
                } else {
                    alert('Réponse non-JSON reçue');
                }
            } else {
                if (response.headers.get('content-type').includes('application/json')) {
                    const errorData = await response.json();
                    alert(`Erreur lors de la connexion: ${errorData.message}`);
                } else {
                    alert('Erreur lors de la connexion');
                }
            }
        } catch (error) {
            console.error('Erreur lors de la connexion:', error);
        }
    });
})
