const express = require('express');
const router = express.Router();
const axios = require('axios');
const urls = require('../utils/urls');
const keys = [
  'AIzaSyAc0FVfz74LKrRtoscnmTSEOZyVC6a6p-o',
  'AIzaSyBtoEtzD2ksmw4JvkOhuMys4cH_-H0VO7w',
];

function request(q, index) {
  return new Promise((resolve, reject) => {
    const options = {
      key: keys[index],
      q: q,
      maxResults: 30,
      part: 'snippet',
      videoCategoryId: 10,
      chart: 'mostpopular',
      type: 'video',
    };

    axios
      .get(urls.baseUrl + '/search', {params: options})
      .then(resp => {
        console.log(resp.status);
        resolve({data: resp.data, status: resp.status});
      })
      .catch(err => {
        console.log(err);
        reject({status: 403, err: err});
      });
  });
}

router.get('/search/:q', (req, res, next) => {
  console.log('search');
  let query = req.params.q;
  let isStatusBad = true;

  let response;
  request(query, 0)
    .then(resp => {
      res.send(resp.data);
    })
    .catch(firsterr => {
      request(query, 1)
        .then(resp => {
          res.send(resp.data);
        })
        .catch(seconderr => {
          console.log(seconderr);
        });
    });
});

module.exports = router;
