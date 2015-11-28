mod = angular.module('mothership')

mod.controller 'PlayerCtrl', ($scope, player, playlist) ->
  'use strict'
  ctrl = this

  $scope.player   = player
  $scope.playlist = playlist
