mod = angular.module('mothership')

mod.controller 'PlaylistCtrl', ($scope, playlist) ->
  'use strict'
  ctrl = this

  $scope.playlist = playlist
