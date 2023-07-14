const express = require("express");
const { Client } = require("youtubei");
const getCircularReplacer = require("../utils/circularDepedencies");
const youtube = new Client();
const router = express.Router();

router.get("/playlist/search/:q", async (req, res) => {
  try {
    const shelves = await youtube.search(`${req.params.q}`, {
      type: "playlist",
    });

    const items = shelves.items.map((item) =>
      JSON.parse(JSON.stringify(item, getCircularReplacer()))
    );

    res.json(items);
  } catch (err) {
    console.log(err);
    res.errored({ message: "Something went wrong" });
  }
});

router.get("/getplaylist/:id", async (req, res) => {
  try {
    const playlist = await youtube.getPlaylist(`${req.params.id}`);

    if (playlist.videos.items.length) {
      const items = playlist.videos.items.map((item) =>
        JSON.parse(JSON.stringify(item, getCircularReplacer()))
      );

      res.json(items);
    }
  } catch (err) {
    console.log(err);
    res.errored({ message: "Something went wrong" });
  }
});

module.exports = router;
