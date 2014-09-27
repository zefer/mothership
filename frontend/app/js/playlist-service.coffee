mod = angular.module("player")

PLAYLIST_STATE_CHANGE = "playlist:state_change"

mod.factory "playlist", ["$rootScope", "mpd", ($rootScope, mpd) ->
  'use strict'
  that = this

  api =
    items: {}
    length: 0
    position: 0

  $rootScope.$on MPD_STATUS, (event, data) ->
    api.position = data.Pos
    api.length = data.playlistlength
    # TODO: only load the playlist when the playlist changed (mpd subsystems)
    mpd.currentPlaylist().then (data) ->
      # TODO: read data into instances of Song or similar?
      api.items = data
      $rootScope.$broadcast PLAYLIST_STATE_CHANGE

  return api
]
