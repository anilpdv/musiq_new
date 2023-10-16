// This route is used to search YouTube for a given query.
// It uses the YouTubei library to make requests to the YouTube API.
// It then converts the shelves object to an array of items and returns it to the client.

const express = require("express");
const router = express.Router();
const { Client } = require("youtubei");
const getCircularReplacer = require("../utils/circularDepedencies");

const youtube = new Client();

// GET   /search/:q
router.get("/search/:q", async (req, res, next) => {
  try {
    const shelves = await youtube.search(`${req.params.q}`, {
      type: "video",
    });

    // Convert the shelves object to an array of items
    const items = shelves.items.map((item) =>
      JSON.parse(JSON.stringify(item, getCircularReplacer()))
    );

    res.json(items);
  } catch (err) {
    console.log(err);
    next(err);
  }
});

module.exports = router;
