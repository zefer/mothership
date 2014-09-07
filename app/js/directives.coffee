mod = angular.module('player')

mod.directive('playerControls', ->
  restrict: 'E'
  templateUrl: 'partials/player-controls.html',
)

mod.directive('playerStatus', ->
  restrict: 'E'
  templateUrl: 'partials/player-status.html',
)
