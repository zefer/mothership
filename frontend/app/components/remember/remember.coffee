mod = angular.module('mothership.remember', [])

mod.factory 'remember', ->
  'use strict'
  that = this

  api =
    get: (key, defaultValue) ->
      api.set(key, defaultValue) unless localStorage.getItem(key)?
      localStorage.getItem(key)

    set: (key, value) -> localStorage.setItem(key, value)
