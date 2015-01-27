mod = angular.module('player')

MPD_STATUS = 'mpd:status'
CONN_STATUS = 'conn:status'

mod.factory 'mpd', ($rootScope, $http, $interval, $q) ->
  ctrl = this
  retrying = null

  @connect = ->
    # Open a websicket and wait for incoming messages.
    ws = new WebSocket("ws://#{location.host}/websocket")
    ws.onopen = ->
      console.log 'Websocket opened'
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
      console.log 'Websocket closed'
      retrying = $interval ->
        console.log 'Websocket reconnecting'
        ctrl.connect()
      , 1000

  api =
    play:       -> $http.get('/play')
    pause:      -> $http.get('/pause')
    previous:   -> $http.get('/previous')
    next:       -> $http.get('/next')
    randomOn:   -> $http.get('/randomOn')
    randomOff:  -> $http.get('/randomOff')

    currentPlaylist: ->
      deferred = $q.defer()
      $http.get('/playlist').success (data) -> deferred.resolve(data)
      deferred.promise

    addToPlaylist: (uri, type, replace, play) ->
      $http.post '/playlist',
        uri: uri
        type: type
        replace: replace
        play: play

    ls: (uri, sort, direction) ->
      deferred = $q.defer()
      $http.get(
        "/files?uri=#{escape(uri)}&sort=#{sort}&direction=#{direction}"
      ).success (data) ->
        deferred.resolve(data)
      deferred.promise

    update: (uri) ->
      $http.put '/library/updated',
        uri: uri

  @connect()

  api
