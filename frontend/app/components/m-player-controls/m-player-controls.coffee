mod = angular.module('mothership.mPlayerControls', [
  'mothership.player'
])

mod.directive 'mPlayerControls', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-player-controls/m-player-controls.html'

  controller: ($scope, player) ->
    vm = this

    $scope.player = player
