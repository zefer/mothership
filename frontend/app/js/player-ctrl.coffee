mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, player, playlist) ->
  'use strict'
  ctrl = this

  $scope.playing  = player.playing
  $scope.playlist = playlist

  $scope.$on PLAYER_STATE_CHANGE, (event) ->
    $scope.playing = player.playing
    $scope.$apply()

  $scope.play     = -> player.play()
  $scope.pause    = -> player.pause()
  $scope.previous = -> player.pause()
  $scope.next     = -> player.next()
  $scope.random   = -> player.random()
)
