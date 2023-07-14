const express = require("express");
const { Client } = require("youtubei");
const getCircularReplacer = require("../utils/circularDepedencies");
const youtube = new Client();
const router = express.Router();

router.get("/playlist/search/:q", async (req, res, next) => {
  try {
    const shelves = await youtube.search(`${req.params.q}`, {
      type: "playlist",
    });

    if (shelves && shelves.items) {
      const items = shelves.items.map((item) =>
        JSON.parse(JSON.stringify(item, getCircularReplacer()))
      );

      res.json(items);
    }
  } catch (err) {
    next(err);
  }
});

router.get("/getplaylist/:id", async (req, res, next) => {
  try {
    const playlist = await youtube.getPlaylist(`${req.params.id}`);

    if (playlist.videos && playlist.videos.items.length) {
      const items = playlist.videos.items.map((item) =>
        JSON.parse(JSON.stringify(item, getCircularReplacer()))
      );

      res.json(items);
    }
  } catch (err) {
    next(err);
  }
});

module.exports = router;
