mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, $interval, $http) ->
  'use strict'
  ctrl = this
  poller = null

  checkPlayerStatus = ->
    console.log('player status poll')
    $http.get('/status').success (data) ->
      $scope.playing =
        now: "#{data.Artist} - #{data.Title}"
        # play, pause or stop
        state: data.state
        error: data.error
        progress: Math.floor((parseFloat(data.elapsed)/parseFloat(data.Time))*100)
        playlistLength: data.playlistlength
        playlistPosition: data.song
        random: data.random == "1"
        quality: ctrl.friendlyQuality(data.audio, data.bitrate)

  ctrl.friendlyQuality = (mpdAudioString, bitrate) ->
    chan = if mpdAudioString.split(':')[2] == '2' then 'Stereo' else 'Mono'
    freq = parseInt(mpdAudioString.split(':')[0]) / 1000 + ' kHz'
    rate = mpdAudioString.split(':')[1] + ' bit'
    bitr = bitrate + ' kbps'
    [chan, rate, freq, bitr].join(', ')

  startMonitoring = ->
    poller = $interval(ctrl.checkPlayerStatus, 5000)

  stopMonitoring = ->
    if angular.isDefined(poller)
      $interval.cancel(poller)
      poller = undefined

  $scope.$on '$destroy', -> $scope.stopMonitoring()

  $scope.play = ->
    $http.get('/play')

  $scope.pause = ->
    $http.get('/pause')

  $scope.previous = ->
    console.log 'previous'
    $http.get('/previous')

  $scope.next = ->
    console.log 'next'
    $http.get('/next')

  $scope.random = ->
    if $scope.playing.random then $http.get('/randomOff') else $http.get('/randomOn')

  startMonitoring()
)
