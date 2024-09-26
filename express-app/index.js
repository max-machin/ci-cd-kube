const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.set('Content-Type', 'text/plain');
  res.send('Hello World');
});

// Route pour tester le 404
app.get('/unknown', (req, res) => {
  res.status(404).send('Not Found');
});

// Route pour tester une erreur 500
app.get('/error', (req, res) => {
  res.status(502).send('Internal Server Error');
});

module.exports = app.listen(8080, () => {
  console.log('Server is running on port 8080');
});
