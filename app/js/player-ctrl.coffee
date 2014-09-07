mod = angular.module('player')

mod.controller('PlayerCtrl', ($scope, $interval) ->
  'use strict'

  monitorPlayer = ->
    checkPlayerStatus = $interval(->
      console.log('player status poll')

      $scope.playing =
        now: 'Joe - ' + Math.random().toString(36).replace(/[^a-z]+/g, '')

    ,1000)

  stopMonitoring = ->
    if angular.isDefined(monitor)
      $interval.cancel(monitor)
      monitor = undefined

  $scope.$on '$destroy', -> $scope.stopMonitoring()

  monitorPlayer()
)
