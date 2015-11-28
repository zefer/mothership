mod = angular.module 'mothership'

mod.directive 'playerControls', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/player-controls/player-controls.html'

  controller: ($scope, player) ->
    vm = this

    $scope.player = player
