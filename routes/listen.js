const ytdl = require("ytdl-core");
const ffmpeg = require("fluent-ffmpeg");
const express = require("express");

const router = express.Router();

const ffmpegPath = require("ffmpeg-static");

ffmpeg.setFfmpegPath(ffmpegPath);

router.get("/listen/:id/:name", (req, res, next) => {
  try {
    let id = req.params.id || "aBt2Djy37tQ";

    let stream = ytdl("https://www.youtube.com/watch?v=" + id, {
      quality: "highestaudio",
      filter: "audioonly",
      highWaterMark: 1 << 25,
    });

    const converter = ffmpeg(stream)
      .audioCodec("libmp3lame")
      .audioBitrate("128k")
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
      res.header("Cache-Control", "public, max-age=3600");
      converter.pipe(res);
    } else {
      res.sendStatus(500);
    }
  } catch (err) {
    next(err);
  }
});

module.exports = router;
