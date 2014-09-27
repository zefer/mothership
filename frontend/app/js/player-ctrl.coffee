mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, player) ->
  'use strict'
  ctrl = this
  $scope.playing = player.playing

  $scope.$on PLAYER_STATE_CHANGE, (event) ->
    $scope.playing = player.playing
    $scope.$apply()

  $scope.play     = -> player.play()
  $scope.pause    = -> player.pause()
  $scope.previous = -> player.pause()
  $scope.next     = -> player.next()
  $scope.random   = -> player.random()
)
