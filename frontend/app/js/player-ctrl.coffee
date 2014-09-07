mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, $interval, $http) ->
  'use strict'
  poller = null

  checkPlayerStatus = ->
    console.log('player status poll')
    $http.get('/status').success (data) ->
      console.log(data)
      $scope.playing =
        now: "#{data.Artist} - #{data.Title}"

  startMonitoring = ->
    poller = $interval(checkPlayerStatus, 1000)

  stopMonitoring = ->
    if angular.isDefined(poller)
      $interval.cancel(poller)
      poller = undefined

  $scope.$on '$destroy', -> $scope.stopMonitoring()

  $scope.play = ->
    console.log 'play/pause'

  $scope.previous = ->
    console.log 'previous'
    $http.get('/previous')

  $scope.next = ->
    console.log 'next'
    $http.get('/next')

  startMonitoring()
)
