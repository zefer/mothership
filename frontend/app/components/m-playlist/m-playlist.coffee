mod = angular.module('mothership.mPlaylist', [
  'mothership.playlist'
])

mod.directive 'mPlaylist', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-playlist/m-playlist.html'

  controller: ($scope, playlist) ->
    vm = this

    $scope.playlist = playlist
