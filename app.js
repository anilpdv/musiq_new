const express = require("express");

// CORS is a node.js package for providing a Connect/Express middleware that can be used to enable CORS with various options.
const cors = require("cors");
// Morgan is a HTTP request logger middleware for Node. js. It simplifies the process of logging requests to your application. You might think of Morgan as a helper that collects logs from your server, such as your request logs.
const morgan = require("morgan");
// swagger is a tool that uses OpenAPI to document APIs
const swaggerUi = require("swagger-ui-express");
// swaggerDocument is a json file that contains all the information about the API
const swaggerDocument = require("./swagger.json");

// importing routes
// search route is used to search for the song
const searchRoute = require("./routes/search");
// listen route is used to get the mp3 file of the song
const listenRoute = require("./routes/listen");
// related route is used to get the related songs
const relatedRoute = require("./routes/relatedSongs.js");
// playlist route is used to get the playlist
const playlistRoute = require("./routes/playlist.js");

const app = express();

// : oas setup
// swaggerUi is the middleware for the swagger
app.use("/docs", swaggerUi.serve, swaggerUi.setup(swaggerDocument));

// morgan is the middleware for the logging
app.use(morgan("combined"));

// cors is the middleware for the cors
app.use(cors());

// server check
// this is the root route
app.get("/", (req, res) => {
  return res.json({
    status: 200,
    routes: {
      searchRoute: "/api/search/:q",
      listenRoute: "/api/listen/:id/:name",
      relatedRoute: "/api/getvideo/:id",
      playlistRoute: "/api/playlist/search/:q",
      playlistRouteById: "/api/getplaylist/:id",
    },
  });
});

// : middle ware
// this is the middleware for the search route
app.use("/api", searchRoute);
// this is the middleware for the listen route
app.use("/api", listenRoute);
// this is the middleware for the related route
app.use("/api", relatedRoute);
// this is the middleware for the playlist route
app.use("/api", playlistRoute);

// : listening to the port
// this is the port number
const port = process.env.PORT || 8080;
// this is the server listening on the port
app.listen(port, () => {
  console.log("server is started and listening on the port " + port);
});
