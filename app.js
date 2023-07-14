// : loading the modules needed
const express = require("express");
const expressOasGenerator = require("express-oas-generator");
const path = require("path");
const cors = require("cors");
const morgan = require("morgan");

// : routes
const searchRoute = require('./routes/search');
const downloadRoute = require('./routes/download');
const lyricsRoute = require('./routes/lyrics');
const listenRoute = require('./routes/listen');
const relatedRoute = require('./routes/relatedSongs.js');
const playlistRoute = require('./routes/playlist.js');
const popularRoute = require('./routes/popular.js');

const app = express();

// : oas setup
app.use(morgan("combined"));
expressOasGenerator.init(app, {});
app.use(cors());

// : middle ware
app.use('/api', searchRoute);
app.use('/api', downloadRoute);
app.use('/api', lyricsRoute);
app.use('/api', listenRoute);
app.use('/api', relatedRoute);
app.use('/api', playlistRoute);
app.use('/api', popularRoute);

// : listening to the port
const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log('server is started and listening on the port ' + port);
});