const ytdl = require("ytdl-core");
const express = require("express");

const router = express.Router();

router.get("/listen/:id/:name", async (req, res, next) => {
  try {
    const { id, name } = req.params;

    // Get the highest quality audio-only stream in mp3 format
    const stream = ytdl("https://www.youtube.com/watch?v=" + id, {
      quality: "highestaudio",
      filter: "audioonly",
      format: "mp3",
      highWaterMark: 1 << 25,
    });

    // Set response headers
    res.header("Content-Disposition", `attachment; filename=${name}`);
    res.header("Cache-Control", "public, max-age=3600");

    // Pipe the stream to the response
    stream.pipe(res);
  } catch (err) {
    next(err);
  }
});

module.exports = router;
