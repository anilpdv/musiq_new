// : loading the modules needed
const express = require("express");
const cors = require("cors");
const morgan = require("morgan");
const swaggerUi = require("swagger-ui-express");
const swaggerDocument = require("./swagger.json");
// : routes
const searchRoute = require("./routes/search");
const listenRoute = require("./routes/listen");
const relatedRoute = require("./routes/relatedSongs.js");
const playlistRoute = require("./routes/playlist.js");

const app = express();

// : oas setup

app.use("/docs", swaggerUi.serve, swaggerUi.setup(swaggerDocument));
app.use(morgan("combined"));
app.use(cors());

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
