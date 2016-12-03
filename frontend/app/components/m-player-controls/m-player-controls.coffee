mod = angular.module('mothership.mPlayerControls', [
  'mothership.player'
])

mod.component 'mPlayerControls',
  bindings: {}
  templateUrl: 'components/m-player-controls/m-player-controls.html'

  controller: (player) ->
    ctrl = this

    ctrl.player = player

    return ctrl
