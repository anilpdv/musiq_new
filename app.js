// : loading the modules needed
const express = require("express");
const expressOasGenerator = require("express-oas-generator");
const path = require("path");
const cors = require("cors");
const morgan = require("morgan");
const compression = require("compression");

// : routes
const searchRoute = require("./routes/search");
const listenRoute = require("./routes/listen");
const relatedRoute = require("./routes/relatedSongs.js");
const playlistRoute = require("./routes/playlist.js");

const app = express();

// : oas setup
app.use(morgan("combined"));
expressOasGenerator.init(app, {});
app.use(cors());
app.use(compression());

// server check
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
app.use("/api", searchRoute);
app.use("/api", listenRoute);
app.use("/api", relatedRoute);
app.use("/api", playlistRoute);

// : listening to the port
const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log("server is started and listening on the port " + port);
});
