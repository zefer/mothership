mod = angular.module('mothership.library', [
  'mothership.mpd'
])

mod.factory 'library', ($q, mpd) ->
  'use strict'

  # Cache a single library path, so the controller can paginate without fetching
  # the data again.
  cache =
    key: null
    items: null

  cacheKey = (uri, sort, direction) ->
    "#{uri}-#{sort}-#{direction}"

  api =
    update: (uri) -> mpd.update(uri)

    ls: (uri, sort, direction) ->
      deferred = $q.defer()
      key = cacheKey(uri, sort, direction)
      if key == cache.key
        # Requesting the same path, return the data from the cache.
        deferred.resolve(cache.items)
      else
        mpd.ls(uri, sort, direction).then (items) ->
          deferred.resolve(items)
          cache.key = key
          cache.items = items
      deferred.promise
