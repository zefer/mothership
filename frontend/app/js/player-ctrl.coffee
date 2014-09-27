mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, $http, player) ->
  'use strict'
  ctrl = this
  $scope.playing = player.playing

  $scope.$on PLAYER_STATE_CHANGE, (event) ->
    $scope.playing = player.playing
    $scope.$apply()

  $scope.play = ->
    $http.get('/play')

  $scope.pause = ->
    $http.get('/pause')

  $scope.previous = ->
    $http.get('/previous')

  $scope.next = ->
    $http.get('/next')

  $scope.random = ->
    if $scope.playing.random then $http.get('/randomOff') else $http.get('/randomOn')
)
