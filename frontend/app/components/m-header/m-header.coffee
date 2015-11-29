mod = angular.module('mothership.mHeader', [
  'mothership.mPlayerControls'
  'mothership.mNavigation'
])

mod.directive 'mHeader', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-header/m-header.html'

  controller: ->
    vm = this
