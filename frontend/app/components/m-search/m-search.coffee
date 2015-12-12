mod = angular.module('mothership.mSearch', [
  'ui.router'

  'mothership.debounce'
])

mod.directive 'mSearch', ->
  restrict: 'E'
  templateUrl: 'components/m-search/m-search.html'
  scope: {}
  controller: 'mSearchController'
  bindToController: true
  controllerAs: 'vm'

mod.controller 'mSearchController', (
  $rootScope, $scope, $state, debounce
) ->
  vm = this

  vm.filter = $state.params.filter

  filter = (filterString) ->
    $state.params.filter = filterString
    $state.go('.', { filter: $state.params.filter }, notify: false)
    $rootScope.$broadcast('search:filter')

  # Don't fire more than once in this many milliseconds, people type fast!
  filter = debounce(filter, 200)

  vm.search = ->
    filter(vm.filter)

  vm.clear = ->
    vm.filter = ''
    vm.search()
