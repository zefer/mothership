mod = angular.module('player')

mod.controller 'PlaylistCtrl', ($scope, playlist) ->
  'use strict'
  ctrl = this

  $scope.playlist = playlist
