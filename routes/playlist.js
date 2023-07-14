const ytlist = require('youtube-playlist');
const express = require('express');
const router = express.Router();

router.get('/playlist/:id', (req, res) => {
  const url = 'https://www.youtube.com/playlist?list=' + req.params.id;

  ytlist(url, ['id', 'name']).then(list => {
    res.json({
      playlist: list.data.playlist,
      list_length: list.data.playlist.length,
    });
  });
});

module.exports = router;
