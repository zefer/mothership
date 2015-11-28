mod = angular.module('mothership')

mod.directive 'playerStatus', ->
  restrict: 'E'
  templateUrl: 'partials/player-status.html'
