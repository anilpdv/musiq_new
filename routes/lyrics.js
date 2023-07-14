const express = require('express');
const axios = require('axios');
var he = require('he');

const convertId = require('../utils/convertId');
const router = express.Router();

router.get('/lyrics/:title', (req, res) => {
  let title;
  let author;
  let query = he.decode(req.params.title);
  console.log(query);
  query = query.split('-');

  console.log(query);
  if (query.length > 1) {
    author = query[0]
      .replace('"', '')
      .split(',')[0]
      .split('&')[0]
      .toLowerCase();

    console.log(author);

    title = query[1]
      .replace(/\"/g, '')
      .replace(/\ /g, '')
      .replace(/\(.*?\)/g, '')
      .replace('&', '')
      .replace('!', '')
      .replace(/\[\]/g, '')
      .split(',')[0];

    console.log(title);
  }

  if (title && author) {
    axios
      .get('https://lyric-api.herokuapp.com/api/find/' + author + '/' + title)
      .then(resp => {
        console.log(resp.data);
        res.json(resp.data);
      })
      .catch(err => {
        console.log(err);
      });
  } else {
    res.json({err: 'lyrics not found'});
  }
});

module.exports = router;
