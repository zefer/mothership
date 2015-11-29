mod = angular.module 'mothership'

mod.directive 'mHeader', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-header/m-header.html'

  controller: ->
    vm = this
