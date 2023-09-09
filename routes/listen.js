// This code is used to convert YouTube videos to mp3 files.
// It uses ytdl-core to get the highest quality audio-only stream from a YouTube video.
// Then it uses ffmpeg to convert the stream to mp3.
// Finally, it pipes the converted stream to the response.

const ytdl = require("ytdl-core");
const ffmpeg = require("fluent-ffmpeg");
const express = require("express");

const router = express.Router();

const ffmpegPath = require("ffmpeg-static");

ffmpeg.setFfmpegPath(ffmpegPath);

router.get("/listen/:id/:name", (req, res, next) => {
  try {
    const { id, name } = req.params;

    // Get the highest quality audio-only stream
    const stream = ytdl("https://www.youtube.com/watch?v=" + id, {
      quality: "highestaudio",
      filter: "audioonly",
      highWaterMark: 1 << 25,
    });

    // Convert the stream to mp3
    const converter = convert(stream);

    if (converter) {
      // Set response headers
      res.header("Content-Disposition", `attachment; filename=${name}`);
      res.header("Cache-Control", "public, max-age=3600");

      // Pipe the converted stream to the response
      converter.pipe(res);
    } else {
      // Return a 500 error if the conversion fails
      res.sendStatus(500);
    }
  } catch (err) {
    next(err);
  }
});

// Convert a stream to mp3
const convert = (stream) => {
  return ffmpeg(stream)
    .audioCodec("libmp3lame")
    .audioQuality(0)
    .format("mp3")
    .on("error", (err) => {
      console.error(err);
      if (!res.headersSent) {
        res.sendStatus(500);
      }
    });
};

module.exports = router;
