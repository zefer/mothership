mod = angular.module("player")

mod.controller("BrowseCtrl", ($scope, $stateParams, $state, library, playlist) ->
  "use strict"
  that = this

  MAX_PER_PAGE = 200

  $scope.library = library

  that.breadcrumbs = (uri) ->
    parts = uri.split("/")
    crumbs = ({
      label: part,
      path: parts[0..i].join("/")
    } for part, i in parts when part != "")
    crumbs.unshift { label: "home", path: "" }
    crumbs

  that.paginate = (items, page) ->
    pages = Math.ceil(items.length / MAX_PER_PAGE)
    page = pages if page > pages
    page = 1 if page < 1
    pos = (page - 1) * MAX_PER_PAGE
    $scope.items = items[pos..pos+MAX_PER_PAGE]
    $scope.pages = (i for i in [1..pages])
    $scope.page = page

  $scope.$on "$stateChangeSuccess", (event, toState, toParams, fromState, fromParams) ->
    toParams.uri ?= "/"
    toParams.page ?= 1
    library.ls(toParams.uri).then (items) ->
      that.paginate(items, parseInt(toParams.page))
      $scope.breadcrumbs = that.breadcrumbs(toParams.uri)

  $scope.showActions = (e) ->
    e.preventDefault()
    e.stopPropagation()
)
