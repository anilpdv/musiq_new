const express = require('express');
const router = express.Router();
const axios = require('axios');
const urls = require('../utils/urls');
const keys = [
  'AIzaSyBc0kmvMa1wbKLQ043RSv-FlhUgQ27pUPk',
  'AIzaSyCUltOMQOkOQkkDnDYjIl3A0QSQ-iUEq8k',
];

const baseUrl = 'https://www.googleapis.com/youtube/v3';

router.get('/popular', (req, res) => {
  // fetch the data : without any search default data
  axios
    .get(baseUrl + '/videos', {
      params: {
        key: keys[0],
        chart: 'mostPopular',
        videoCategoryId: 10,
        part: 'snippet',
        maxResults: 20,
      },
    })
    .then(function(resp) {
      res.json(resp.data);
    })
    .catch(function(err) {
      res.json(err.response.data);
    });
});

module.exports = router;
