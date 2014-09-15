mod = angular.module("player")

MPD_STATUS = "mpd:status"

mod.factory "mpdService", ["$rootScope", ($rootScope) ->

  # Open a websicket and wait for incoming messages.
  ws = new WebSocket("ws://#{location.host}/websocket")
  ws.onopen = ->
    console.log "Websocket opened."

  # The only type of message is mpd status JSON, for now.
  ws.onmessage = (message) ->
    $rootScope.$broadcast MPD_STATUS, JSON.parse(message.data)

]
