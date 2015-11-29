mod = angular.module 'mothership'

mod.directive 'mPlaying', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-playing/m-playing.html'

  controller: ($scope, player) ->
    vm = this

    $scope.player = player
