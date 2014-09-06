(->
  'use strict'

  angular.module('player', [])

    .controller('PlayerController', ($scope, $interval) ->

      monitorPlayer = ->
        checkPlayerStatus = $interval(->
          console.log('player status poll')

          $scope.playing =
            artist: 'joe'
            song: Math.random().toString(36).replace(/[^a-z]+/g, '')

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
