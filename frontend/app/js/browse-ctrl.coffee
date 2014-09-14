mod = angular.module("player")

mod.controller("BrowseCtrl", ($scope, $stateParams, $state, $http) ->
  "use strict"
  ctrl = this
  $scope.uri = ""

  $scope.ls = (uri) ->
    $http.get("/files?uri=#{escape(uri)}").success (data) ->
      $scope.items = data
      $scope.uri = uri

  $scope.$on "$stateChangeSuccess", (event, toState, toParams, fromState, fromParams) ->
    toParams.uri ?= "/"
    $scope.ls toParams.uri
)
