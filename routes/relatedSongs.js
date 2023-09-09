// This code gets a video from YouTube.
// It then returns the related videos.
// The video ID is passed as a parameter to the function.
// The function returns a list of related videos.

const express = require("express");
const { Client } = require("youtubei");

const router = express.Router();
const getCircularReplacer = require("../utils/circularDepedencies");
const youtube = new Client();

const getVideo = async (id) => {
  // Get a video from YouTube.
  const video = await youtube.getVideo(`${req.params.id}`);

  // Get the related videos.
  if (video && video.related) {
    // Get the first page of related videos.
    await video.related.next(0);

    // Convert the items to a JSON string.
    const items = video.related.items.map((item) =>
      JSON.parse(JSON.stringify(item, getCircularReplacer()))
    );

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

module.exports = router;
