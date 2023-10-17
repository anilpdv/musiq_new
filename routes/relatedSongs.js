// This code gets a video from YouTube.
// It then returns the related videos.
// The video ID is passed as a parameter to the function.
// The function returns a list of related videos.

const express = require("express");
const { Client } = require("youtubei");
const ytdl = require("ytdl-core");
const router = express.Router();
const getCircularReplacer = require("../utils/circularDepedencies");
const youtube = new Client();

const getVideo = async (id) => {
  // Get a video from YouTube.
  const video = await youtube.getVideo(id, { type: "video" });

  // Get the related videos.
  if (video && video.related) {
    // Get the first page of related videos.
    await video.related.next(0);

    // Sort the items by viewCount.
    const items = video.related.items
      .sort((a, b) => b.viewCount - a.viewCount)
      .map((item) => JSON.parse(JSON.stringify(item, getCircularReplacer())));

    return items;
  }
};

// Get a video from YouTube.
router.get("/getvideo/:id", async (req, res, next) => {
  try {
    // Get the related videos.
    const items = await getVideo(req.params.id);

    // Return the related videos.
    res.json(items);
  } catch (err) {
    // Log any errors.
    console.error(err);

    // Return an error message.
    res.status(500).send("Something went wrong");
  }
});

router.get("/related/:id", async (req, res, next) => {
  try {
    let id = req.params.id;
    const info = await ytdl.getInfo(id);
    const responseJson = {
      videoDetails: info.videoDetails,
      relatedSongs: info.related_videos,
    };
    res.json(responseJson);
  } catch (err) {
    console.error(err);
    res.status(500).send("Something went wrong");
  }
});

module.exports = router;
