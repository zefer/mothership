mod = angular.module('mothership')

mod.directive 'sortBy', ->
  restrict: 'E'
  templateUrl: 'components/sort-by/sort-by.html'
  scope: {}

  controller: ->
    vm = this
