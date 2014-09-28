mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, player, playlist) ->
  'use strict'
  ctrl = this

  $scope.player   = player
  $scope.playlist = playlist
)
