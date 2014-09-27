mod = angular.module("player")

PLAYER_STATE_CHANGE = "player:state_change"

mod.factory "player", ["$rootScope", "mpdService", ($rootScope, mpdService) ->
  'use strict'
  that = this
  # The public methods/data we expose
  api =
    playing: {}

  $rootScope.$on MPD_STATUS, (event, data) ->
    [now, sub] = that.nowPlaying(data)
    api.playing =
      now: now
      sub: sub
      # play, pause or stop
      state: data.state
      error: data.error
      progress: Math.floor((parseFloat(data.elapsed)/parseFloat(data.Time))*100)
      playlistLength: data.playlistlength
      playlistPosition: parseInt(data.song||-1) + 1
      random: data.random == "1"
      quality: that.friendlyQuality(data.audio, data.bitrate)
    $rootScope.$broadcast PLAYER_STATE_CHANGE

  $rootScope.$on CONN_STATUS, (event, connected) ->
    api.playing.error = if connected then "" else "Connection lost"
    $rootScope.$broadcast PLAYER_STATE_CHANGE

  that.nowPlaying = (data) ->
    if data.Artist && data.Title
      now = "#{data.Artist} - #{data.Title}"
      sub = data.Album
    else if data.Name
      now = data.Name
      sub = ""
    else
      parts = data.file.split("/")
      now = parts[parts.length-1]
      sub = parts[0..parts.length-2].join("/")
    [now, sub]

  that.friendlyQuality = (mpdAudioString, bitrate) ->
    return unless mpdAudioString
    chan = if mpdAudioString.split(':')[2] == '2' then 'Stereo' else 'Mono'
    freq = parseInt(mpdAudioString.split(':')[0]) / 1000 + ' kHz'
    rate = mpdAudioString.split(':')[1] + ' bit'
    bitr = bitrate + ' kbps'
    [chan, rate, freq, bitr].join(', ')

  return api
]
