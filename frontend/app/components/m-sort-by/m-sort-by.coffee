mod = angular.module('mothership')

mod.directive 'mSortBy', ->
  restrict: 'E'
  templateUrl: 'components/m-sort-by/m-sort-by.html'
  scope: {}

  controller: ->
    vm = this
