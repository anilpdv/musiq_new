const express = require('express');
const axios = require('axios');
const ytdl = require('youtube-dl');
const fs = require('fs');

const router = express.Router();

router.get('/download/:id/song/:name', (req, res) => {
  const video = ytdl('https://www.youtube.com/watch?v=' + req.params.id, [
    '--format=140',
  ]);
  res.setHeader('Content-Type', 'audio/mpeg');
  res.setHeader(
    'Content-Disposition',
    'attachment; filename=' + req.params.name + '.m4a',
  );
  video.on('error', err => console.log(err));
  video.pipe(res);
});

module.exports = router;
