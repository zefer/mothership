mod = angular.module 'mothership'

mod.directive 'playerStatus', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/player-status/player-status.html'

  controller: ($scope, player) ->
    vm = this

    $scope.player = player
