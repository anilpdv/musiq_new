const express = require("express");
const router = express.Router();
const { Client } = require("youtubei");
const getCircularReplacer = require("../utils/circularDepedencies");

const youtube = new Client();

router.get("/search/:q", async (req, res, next) => {
  try {
    const shelves = await youtube.search(`${req.params.q}`, {
      type: "video",
    });

    const items = shelves.items.map((item) =>
      JSON.parse(JSON.stringify(item, getCircularReplacer()))
    );

    res.json(items);
  } catch (err) {
    next(err);
  }
});

module.exports = router;
