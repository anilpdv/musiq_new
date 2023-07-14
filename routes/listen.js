const ytdl = require("ytdl-core");
const ffmpeg = require("fluent-ffmpeg");
const express = require("express");

const router = express.Router();

const ffmpegPath = require("ffmpeg-static");

ffmpeg.setFfmpegPath(ffmpegPath);

router.get("/listen/:id/:name", (req, res) => {
  let id = req.params.id || "aBt2Djy37tQ";

  let stream = ytdl("https://www.youtube.com/watch?v=" + id, {
    quality: "highestaudio",
  });

  const converter = ffmpeg(stream)
    .audioBitrate(128)
    .format("mp3")
    .on("error", (err) => {
      console.error(err);
      if (!res.headersSent) {
        res.sendStatus(500);
      }
    });

  if (converter) {
    res.header(
      "Content-Disposition",
      `attachment; filename=${req.params.name}`
    );
    converter.pipe(res);
  } else {
    res.sendStatus(500);
  }
});

module.exports = router;
