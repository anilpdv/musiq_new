const express = require("express");
const { Client } = require("youtubei");

const router = express.Router();
const getCircularReplacer = require("../utils/circularDepedencies");
const youtube = new Client();

router.get("/getvideo/:id", async (req, res, next) => {
  try {
    const video = await youtube.getVideo(`${req.params.id}`);

    if (video && video.related) {
      await video.related.next(0);
      const items = video.related.items.map((item) =>
        JSON.parse(JSON.stringify(item, getCircularReplacer()))
      );
      res.json(items);
    }
  } catch (err) {
    next(err);
  }
});

module.exports = router;
