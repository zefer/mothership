mod = angular.module("player")

mod.controller("BrowseCtrl", ($scope, $stateParams, $state, $http) ->
  "use strict"
  ctrl = this
  $scope.uri = ""

  $scope.ls = (uri) ->
    $http.get("/files?uri=#{escape(uri)}").success (data) ->
      $scope.items = data
      $scope.uri = uri
      $scope.breadcrumbs = ctrl.breadcrumbs(uri)

  ctrl.breadcrumbs = (uri) ->
    parts = uri.split("/")
    { label: part, path: parts[0..i].join("/") } for part, i in parts

  ctrl.add = (uri, replace, play) ->
    $http.post '/playlist',
      uri: uri
      replace: replace
      play: play

  $scope.$on "$stateChangeSuccess", (event, toState, toParams, fromState, fromParams) ->
    toParams.uri ?= "/"
    $scope.ls toParams.uri

  $scope.showActions = (e) ->
    e.preventDefault()
    e.stopPropagation()

  $scope.add = (uri) ->
    ctrl.add(uri, false, false)

  $scope.addPlay = (uri) ->
    ctrl.add(uri, false, true)

  $scope.addReplacePlay = (uri) ->
    ctrl.add(uri, true, true)

  $scope.update = (uri) ->
    console.log "update", uri

)
