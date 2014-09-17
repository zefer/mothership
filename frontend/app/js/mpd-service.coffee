mod = angular.module("player")

MPD_STATUS = "mpd:status"
CONN_STATUS = "conn:status"

mod.factory "mpdService", ["$rootScope", "$interval", ($rootScope, $interval) ->
  ctrl = this
  retrying = null

  @connect = ->
    # Open a websicket and wait for incoming messages.
    ws = new WebSocket("ws://#{location.host}/websocket")
    ws.onopen = ->
      console.log "Websocket opened"
      retrying && $interval.cancel(retrying)
      retrying = null
      $rootScope.$broadcast CONN_STATUS, true

    # The only type of message is mpd status JSON, for now.
    ws.onmessage = (message) ->
      $rootScope.$broadcast MPD_STATUS, JSON.parse(message.data)

    # The only type of message is mpd status JSON, for now.
    ws.onclose = ->
      return if retrying
      $rootScope.$broadcast CONN_STATUS, false
      console.log "Websocket closed"
      retrying = $interval ->
        console.log "Websocket reconnecting"
        ctrl.connect()
      , 1000

  @connect()
]
