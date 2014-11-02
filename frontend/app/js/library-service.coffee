mod = angular.module("player")

mod.factory "library", ["$q", "mpd", ($q, mpd) ->
  'use strict'

  # Cache a single library path, so the controller can paginate without fetching
  # the data again.
  cache =
    uri: null
    items: null

  api =
    update: (uri) -> mpd.update(uri)

    ls: (uri) ->
      deferred = $q.defer()
      if uri == cache.uri
        # Requesting the same path, return the data from the cache.
        deferred.resolve(cache.items)
      else
        mpd.ls(uri).then (items) ->
          deferred.resolve(items)
          cache.uri = uri
          cache.items = items
      deferred.promise
]
