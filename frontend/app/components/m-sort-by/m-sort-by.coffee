mod = angular.module('mothership.mSortBy', [])

mod.directive 'mSortBy', ->
  restrict: 'E'
  templateUrl: 'components/m-sort-by/m-sort-by.html'
  scope: {}

  controller: ->
    vm = this
