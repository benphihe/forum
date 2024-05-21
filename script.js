document.addEventListener('DOMContentLoaded', function() {
    // Popup Inscription
    document.getElementById('openInscriptionPopup').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreInscription').style.display = 'block';
        document.getElementById('popupOverlay').style.display = 'block';
    });

    document.getElementById('closeInscriptionPopup').addEventListener('click', function() {
        document.getElementById('fenetreInscription').style.display = 'none';
        document.getElementById('popupOverlay').style.display = 'none';
    });

    // Popup Connexion
    document.getElementById('openConnexionPopup').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreConnexion').style.display = 'block';
        document.getElementById('popupOverlay').style.display = 'block';
    });

    document.getElementById('closeConnexionPopup').addEventListener('click', function() {
        document.getElementById('fenetreConnexion').style.display = 'none';
        document.getElementById('popupOverlay').style.display = 'none';
    });

    // Popup Inscription depuis Popup Connexion
    document.getElementById('openInscriptionPopupFromConnexion').addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('fenetreConnexion').style.display = 'none'; // Fermer le popup de connexion
        document.getElementById('fenetreInscription').style.display = 'block'; // Ouvrir le popup d'inscription
        document.getElementById('popupOverlay').style.display = 'block'; // Afficher l'overlay
    });
});
