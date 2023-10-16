// This code is used to convert YouTube videos to mp3 files.
// It uses ytdl-core to get the highest quality audio-only stream from a YouTube video.
// Then it uses ffmpeg to convert the stream to mp3.
// Finally, it pipes the converted stream to the response.

const ytdl = require("ytdl-core");
const express = require("express");
const router = express.Router();
const ffmpeg = require("fluent-ffmpeg");
const ffmpegPath = require("ffmpeg-static");
const cp = require("child_process");
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
    const converter = convert(stream, res);

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

router.get("/watch/:id/:name", async (req, res) => {
  const { id, name } = req.params;
  let url = "https://www.youtube.com/watch?v=" + id;

  res.header("Content-Disposition", `attachment; filename=${name}`);

  let video = ytdl(url, { filter: "videoonly" });
  let audio = ytdl(url, { filter: "audioonly", highWaterMark: 1 << 25 });

  const ffmpegProcess = cp.spawn(
    ffmpegPath,
    [
      "-i",
      "pipe:3",
      "-i",
      "pipe:4",
      "-map",
      "0:v",
      "-map",
      "1:a",
      "-c:v",
      "copy",
      "-c:a",
      "libmp3lame",
      "-crf",
      "27",
      "-preset",
      "veryfast",
      "-movflags",
      "frag_keyframe+empty_moov",
      "-f",
      "mp4",
      "-loglevel",
      "error",
      "-",
    ],
    {
      stdio: ["pipe", "pipe", "pipe", "pipe", "pipe"],
    }
  );

  video.pipe(ffmpegProcess.stdio[3]);
  audio.pipe(ffmpegProcess.stdio[4]);
  ffmpegProcess.stdio[1].pipe(res);

  let ffmpegLogs = "";

  ffmpegProcess.stdio[2].on("data", (chunk) => {
    ffmpegLogs += chunk.toString();
  });

  ffmpegProcess.on("exit", (exitCode) => {
    if (exitCode === 1) {
      console.error(ffmpegLogs);
    }
  });
});
// Convert a stream to mp3
const convert = (stream, res) => {
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
