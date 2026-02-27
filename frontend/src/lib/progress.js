// Client-side progress tracking: interpolates elapsed time between server updates.
// The server only pushes status on MPD events (play/pause/track change), so we use
// a wall-clock timer to keep the progress bar and elapsed time ticking each second.

import { playerStatus, elapsedTime } from './stores.js';

let timer = null;
let trackEndTime = null;
let trackDuration = 0;

function startTimer() {
  stopTimer();
  timer = setInterval(() => {
    const remaining = (trackEndTime - Date.now()) / 1000;
    const elapsed = Math.max(0, Math.min(trackDuration, trackDuration - remaining));
    elapsedTime.set(elapsed);
  }, 1000);
}

function stopTimer() {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
}

// Sync timer with server status updates.
playerStatus.subscribe(($status) => {
  const state = $status.state || 'stop';
  const elapsed = parseFloat($status.elapsed) || 0;
  const duration = parseFloat($status.duration) || parseFloat($status.Time) || 0;

  elapsedTime.set(elapsed);
  trackDuration = duration;

  if (state === 'play' && duration > 0) {
    trackEndTime = Date.now() + (duration - elapsed) * 1000;
    startTimer();
  } else {
    stopTimer();
  }
});
