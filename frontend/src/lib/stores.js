import { writable, derived } from 'svelte/store';

// WebSocket connection status.
export const connected = writable(false);

// Raw MPD status data from WebSocket.
export const playerStatus = writable({});

// Current playlist items (fetched via HTTP).
export const playlist = writable([]);

// Client-side elapsed time, updated by progress.js timer.
export const elapsedTime = writable(0);

// Derived player state.
export const player = derived([playerStatus, elapsedTime], ([$status, $elapsed]) => {
  const state = $status.state || 'stop';
  const error = $status.error || '';
  const randomOn = $status.random === '1';

  // Now playing display.
  let now = '';
  let sub = '';
  if ($status.Artist && $status.Title) {
    now = `${$status.Artist} - ${$status.Title}`;
    sub = $status.Album || '';
  } else if ($status.Name) {
    now = $status.Name;
  } else if ($status.file) {
    const parts = $status.file.split('/');
    now = parts[parts.length - 1];
    sub = parts.slice(0, -1).join('/');
  }

  // Audio quality.
  let quality = '';
  if ($status.audio) {
    quality = formatQuality($status.audio, $status.bitrate);
  }

  // Progress.
  const duration = parseFloat($status.duration) || parseFloat($status.Time) || 0;
  const percentage = duration > 0 ? Math.round(($elapsed / duration) * 100) : 0;

  return {
    state,
    error,
    randomOn,
    now,
    sub,
    quality,
    elapsed: formatTime($elapsed),
    total: formatTime(duration),
    percentage,
    // Raw position in playlist (1-indexed, from MPD status).
    position: $status.Pos ? parseInt($status.Pos, 10) + 1 : 0,
    playlistLength: parseInt($status.playlistlength, 10) || 0,
  };
});

function formatTime(seconds) {
  if (!seconds || seconds <= 0) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatQuality(audio, bitrate) {
  const parts = audio.split(':');
  if (parts.length < 3) return '';

  const sampleRate = parseInt(parts[0], 10);
  const bits = parts[1];
  const channels = parseInt(parts[2], 10);

  const channelStr = channels === 2 ? 'Stereo' : channels === 1 ? 'Mono' : `${channels}ch`;
  const sampleStr = sampleRate >= 1000 ? `${(sampleRate / 1000).toFixed(1)} kHz` : `${sampleRate} Hz`;

  let result = `${channelStr}, ${bits} bit, ${sampleStr}`;
  if (bitrate && bitrate !== '0') {
    result += `, ${bitrate} kbps`;
  }
  return result;
}
