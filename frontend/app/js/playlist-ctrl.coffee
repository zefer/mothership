mod = angular.module("player")

mod.controller("PlaylistCtrl", ($rootScope, $scope, $stateParams, $state, $http) ->
  "use strict"
  ctrl = this

  ctrl.update = ->
    $http.get("/playlist").success (data) ->
      $scope.items = data

  $scope.$on "$stateChangeSuccess", (event, toState, toParams, fromState, fromParams) ->
    ctrl.update()

)
