mod = angular.module 'mothership'

mod.directive 'mPlayerStatus', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-player-status/m-player-status.html'

  controller: ($scope, player) ->
    vm = this

    $scope.player = player
