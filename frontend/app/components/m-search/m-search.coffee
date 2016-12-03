mod = angular.module('mothership.mSearch', [
  'ui.router'

  'mothership.debounce'
])

mod.component 'mSearch',
  templateUrl: 'components/m-search/m-search.html'

  bindings: {}

  controller: ($rootScope, $state, debounce) ->
    ctrl = this

    ctrl.filter = $state.params.filter

    filter = (filterString) ->
      $state.params.filter = filterString
      $state.go('.', { filter: $state.params.filter }, notify: false)
      $rootScope.$broadcast('search:filter')

    # Don't fire more than once in this many milliseconds, people type fast!
    filter = debounce(filter, 200)

    ctrl.search = ->
      filter(ctrl.filter)

    ctrl.clear = ->
      ctrl.filter = ''
      ctrl.search()

    return ctrl
