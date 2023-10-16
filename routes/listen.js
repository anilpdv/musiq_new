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
  const url = "https://www.youtube.com/watch?v=" + id;
  const info = await ytdl.getInfo(id);

  // Ensure the client supports range requests
  const range = req.headers.range;
  if (!range) {
    res.status(416).send("Range request required");
    return;
  }

  const video = await ytdl(url, {
    quality: "highestvideo",
    filter: "videoonly",
    highWaterMark: 1 << 25,
  });

  const audioFormats = info.formats
    .filter((f) => f.mimeType.includes("audio"))
    .sort((a, b) => b.audioBitrate - a.audioBitrate);

  const englishAudio = audioFormats.find((f) => {
    return f.audioTrack && f.audioTrack.id.startsWith("en");
  });

  const audio = await ytdl(url, {
    format: englishAudio || audioFormats[0],
  });

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
      "aac",
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

  const totalSize = video._readableState.length + audio._readableState.length;

  const head = {
    "Content-Length": totalSize,
    "Content-Type": "video/mp4",
  };

  res.writeHead(200, head);

  ffmpegProcess.stdio[1].pipe(res);

  video.pipe(ffmpegProcess.stdio[3]);
  audio.pipe(ffmpegProcess.stdio[4]);
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

// router for geting info json of video id
router.get("/info/:id", async (req, res) => {
  try {
    const { id } = req.params;
    const info = await ytdl.getInfo(id);
    res.json(info);
  } catch (err) {
    console.error(err);
    res.sendStatus(500);
  }
});

module.exports = router;
