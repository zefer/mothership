mod = angular.module("player")

mod.controller("BrowseCtrl", ($scope, $stateParams, $state, library, playlist) ->
  "use strict"
  that = this

  $scope.library = library

  that.breadcrumbs = (uri) ->
    parts = uri.split("/")
    crumbs = ({
      label: part,
      path: parts[0..i].join("/")
    } for part, i in parts when part != "")
    crumbs.unshift { label: "home", path: "" }
    crumbs

  $scope.$on "$stateChangeSuccess", (event, toState, toParams, fromState, fromParams) ->
    toParams.uri ?= "/"
    library.ls(toParams.uri).then (data) ->
      $scope.items = data
      $scope.breadcrumbs = that.breadcrumbs(toParams.uri)

  $scope.showActions = (e) ->
    e.preventDefault()
    e.stopPropagation()
)
