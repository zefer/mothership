mod = angular.module("player")

mod.factory "library", ["mpd", (mpd) ->
  'use strict'

  api =
    ls:     (uri)  -> mpd.ls(uri)
    update: (uri)  -> mpd.update(uri)
]
