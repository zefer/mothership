(->
  'use strict'

  angular.module('player', [])

    .controller('PlayerController', ($scope, $interval) ->

      monitorPlayer = ->
        checkPlayerStatus = $interval(->
          console.log('player status poll')
        ,1000)

      stopMonitoring = ->
        if angular.isDefined(monitor)
          $interval.cancel(monitor)
          monitor = undefined

      $scope.$on('$destroy', ->
        $scope.stopMonitoring()
      )

      monitorPlayer()
    )

    .directive('playerControls', ->
      restrict: 'E'
      templateUrl: 'partials/player-controls.html',
    )

    .directive('playerStatus', ->
      restrict: 'E'
      templateUrl: 'partials/player-status.html',
    )

)()
