// Import the express library
const express = require("express");
// Import the youtube client
const { Client } = require("youtubei");
// Import the circular dependencies fix
const getCircularReplacer = require("../utils/circularDepedencies");

// Create a new youtube client
const youtube = new Client();
// Create a new express router
const router = express.Router();

// Create a route to get playlists that match a query
router.get("/playlist/search/:q", async (req, res, next) => {
  try {
    // Search for playlists that match a query
    const shelves = await youtube.search(`${req.params.q}`, {
      type: "playlist",
    });

    // Check if the search returned any playlists
    if (shelves && shelves.items) {
      // Convert the playlists to JSON
      const items = shelves.items.map((item) =>
        JSON.parse(JSON.stringify(item, getCircularReplacer()))
      );

      // Return the playlists
      res.json(items);
    }
  } catch (err) {
    next(err);
  }
});

// Create a route to get the items in a playlist
router.get("/getplaylist/:id", async (req, res, next) => {
  try {
    // Get the playlist with the specified ID
    const playlist = await youtube.getPlaylist(`${req.params.id}`);

    // Get the items in the playlist
    const items = await getPlaylistItems(playlist);

    // Return the items in the playlist
    res.json(items);
  } catch (err) {
    next(err);
  }
});

// Get the items in a playlist
const getPlaylistItems = async (playlist) => {
  // Check if there are any items in the playlist
  if (playlist.videos && playlist.videos.items.length) {
    // Convert the items to JSON
    const items = playlist.videos.items.map((item) =>
      JSON.parse(JSON.stringify(item, getCircularReplacer()))
    );

    // Return the items
    return items;
  } else {
    // Throw an error if no items were found
    throw new Error("No items found");
  }
};

// Export the router
module.exports = router;
