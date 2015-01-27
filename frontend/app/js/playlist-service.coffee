mod = angular.module("player")

PLAYLIST_STATE_CHANGE = "playlist:state_change"

mod.factory "playlist", ($rootScope, mpd) ->
  'use strict'
  that = this

  api =
    items: {}
    length: 0
    position: 0

    add:            (uri, type) -> mpd.addToPlaylist(uri, type, false, false)
    addPlay:        (uri, type) -> mpd.addToPlaylist(uri, type, false, true)
    addReplacePlay: (uri, type) -> mpd.addToPlaylist(uri, type, true,  true)

  $rootScope.$on MPD_STATUS, (event, data) ->
    api.position = parseInt(data.Pos||-1)+1
    api.length = data.playlistlength
    # TODO: only load the playlist when the playlist changed (mpd subsystems)
    mpd.currentPlaylist().then (data) ->
      # TODO: read data into instances of Song or similar?
      api.items = data
      $rootScope.$broadcast PLAYLIST_STATE_CHANGE

  return api
