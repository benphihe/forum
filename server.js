const express = require('express');
const bodyParser = require('body-parser');
const sqlite3 = require('sqlite3').verbose();
const path = require('path');

const app = express();
const port = 5000;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.use(express.static(path.join(__dirname)));

const dbname = 'ma_base_de_donnees.db';
let db = new sqlite3.Database(dbname, (err) => {
    if (err) {
        console.error(err.message);
        throw err;
    }
    console.log('Connected to the SQLite database.');

    db.run(`CREATE TABLE IF NOT EXISTS utilisateurs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        fullname TEXT,
        email TEXT UNIQUE,
        mdp TEXT
    )`, (err) => {
        if (err) {
            console.log('Table already exists.');
        }
    });
});

app.post('/inscription', (req, res) => {
    const { fullname, email, mdp } = req.body;
    const sql = 'INSERT INTO utilisateurs (fullname, email, mdp) VALUES (?, ?, ?)';
    db.run(sql, [fullname, email, mdp], function(err) {
        if (err) {
            return res.status(400).json({ error: err.message });
        }
        res.json({ message: 'Inscription réussie', id: this.lastID });
    });
});

app.post('/connexion', (req, res) => {
    const { email, mdp } = req.body;
    const sql = 'SELECT * FROM utilisateurs WHERE email = ? AND mdp = ?';
    db.get(sql, [email, mdp], (err, row) => {
        if (err) {
            return res.status(400).json({ error: err.message });
        }
        if (row) {
            res.json({ message: 'Connexion réussie', user: row });
        } else
        {
          res.status(401).json({ message: 'Email ou mot de passe incorrect' });
      }
  });
});

app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'Accueil.html'));
});

app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});





