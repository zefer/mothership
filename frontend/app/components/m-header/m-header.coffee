mod = angular.module('mothership.mHeader', [
  'mothership.mPlayerControls'
  'mothership.mNavigation'
])

mod.component 'mHeader',
  bindings: {}
  templateUrl: 'components/m-header/m-header.html'

  controller: ->
    ctrl = this

    return ctrl
