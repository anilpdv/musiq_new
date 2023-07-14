const ytdl = require('ytdl-core');
const ffmpeg = require('fluent-ffmpeg');
const express = require('express');
const router = express.Router();

// ffmpeg.setFfmpegPath("path/to/your/ffmpeg"); // Set this to the path where you have FFmpeg installed.

router.get('/listen/:id', (req, res) => {
    let stream = ytdl('https://www.youtube.com/watch?v=' + req.params.id, {
        quality: 'highestaudio',
    });

    res.header('Content-Disposition', `attachment; filename="audio.mp3"`);
    ffmpeg(stream)
        .audioBitrate(128)
        .format('mp3')
        .on('error', (err) => {
            console.error(err);
            res.sendStatus(500);
        })
        .pipe(res);
});

module.exports = router;
