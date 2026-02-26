// WebSocket + HTTP API client for MPD backend.

import { connected, playerStatus, playlist } from './stores.js';

let ws = null;
let retryInterval = null;

function wsUrl() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${location.host}/websocket`;
}

export function connectWebSocket() {
  if (ws) return;
  ws = new WebSocket(wsUrl());

  ws.onopen = () => {
    connected.set(true);
    if (retryInterval) {
      clearInterval(retryInterval);
      retryInterval = null;
    }
  };

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    playerStatus.set(data);
    // Refresh playlist on every status update.
    fetchPlaylist();
  };

  ws.onclose = () => {
    connected.set(false);
    ws = null;
    if (!retryInterval) {
      retryInterval = setInterval(connectWebSocket, 1000);
    }
  };

  ws.onerror = () => {
    ws?.close();
  };
}

// Player commands.
export const play = () => fetch('/play');
export const pause = () => fetch('/pause');
export const next = () => fetch('/next');
export const previous = () => fetch('/previous');
export const randomOn = () => fetch('/randomOn');
export const randomOff = () => fetch('/randomOff');

// Browse the music library.
export async function ls(uri, sort, direction, filter) {
  const params = new URLSearchParams({ uri });
  if (sort) params.set('sort', sort);
  if (direction) params.set('direction', direction);
  if (filter) params.set('filter', filter);
  const res = await fetch(`/files?${params}`);
  return res.json();
}

// Get the current playlist.
async function fetchPlaylist() {
  try {
    const res = await fetch('/playlist');
    const items = await res.json();
    playlist.set(items);
  } catch {
    // Ignore fetch errors during reconnect.
  }
}

// Add to playlist.
export function addToPlaylist(uri, type, replace = false, playAfter = false) {
  return fetch('/playlist', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ uri, type, replace, play: playAfter }),
  });
}

// Update (rescan) a library folder.
export function updateLibrary(uri) {
  return fetch('/library/updated', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ uri }),
  });
}
