const express = require('express');
const router = express.Router();
const axios = require('axios');
const urls = require('../utils/urls');
const keys = [
  'AIzaSyBc0kmvMa1wbKLQ043RSv-FlhUgQ27pUPk',
  'AIzaSyCUltOMQOkOQkkDnDYjIl3A0QSQ-iUEq8k',
];

function request(id, index) {
  return new Promise((resolve, reject) => {
    const options = {
      key: keys[index],
      maxResults: 20,
      part: 'snippet',
      videoCategoryId: 10,
      chart: 'mostpopular',
      type: 'video',
      relatedToVideoId: id,
    };

    axios
      .get(urls.baseUrl + '/search', {params: options})
      .then(resp => {
        resolve(resp.data);
      })
      .catch(err => {
        reject(err);
      });
  });
}

router.get('/related/:id', (req, res, next) => {
  console.log('related', req.params.id);
  let query = req.params.id;

  request(query, 0)
    .then(data => {
      console.log('first request success');
      res.send(data);
    })
    .catch(firsterr => {
      request(query, 1)
        .then(data => {
          console.log('second request succes');
          res.send(data);
        })
        .catch(seconderr => {
          console.log(seconderr);
        });
    });
});

module.exports = router;
