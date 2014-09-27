mod = angular.module("player")

PLAYER_STATE_CHANGE = "player:state_change"

mod.factory "player", ["$rootScope", "$http", "mpd", ($rootScope, $http, mpd) ->
  'use strict'
  that = this

  api =
    play:    -> mpd.play()
    pause:   -> mpd.pause()
    previous:-> mpd.previous()
    next:    -> mpd.next()
    random:  -> if api.randomOn then mpd.randomOff() else mpd.randomOn()

  $rootScope.$on MPD_STATUS, (event, data) ->
    [now, sub] = that.nowPlaying(data)
    # TODO: tidy up this exposed data, reuse a Song class to represent playing?
    angular.extend api,
      now: now
      sub: sub
      # play, pause or stop
      state: data.state
      error: data.error
      progress: Math.floor((parseFloat(data.elapsed)/parseFloat(data.Time))*100)
      randomOn: data.random == "1"
      quality: that.friendlyQuality(data.audio, data.bitrate)
    $rootScope.$broadcast PLAYER_STATE_CHANGE

  $rootScope.$on CONN_STATUS, (event, connected) ->
    api.error = if connected then "" else "Connection lost"
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
