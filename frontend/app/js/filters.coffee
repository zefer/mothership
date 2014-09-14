mod = angular.module("player")

mod.filter "basename", ->
  (path) ->
    parts = (path || "").split("/")
    parts[parts.length-1] || ""
