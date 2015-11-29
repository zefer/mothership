mod = angular.module('mothership.mPagination', [])

mod.directive 'mPagination', ->
  restrict: 'E'
  scope:
    pages: '='
    page: '='
  templateUrl: 'components/m-pagination/m-pagination.html'

  controller: ($scope) ->
    vm = this
