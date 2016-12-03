mod = angular.module('mothership.mBrowseActions', [
  'mothership.playlist'
  'mothership.library'
])

mod.component 'mBrowseActions',
  templateUrl: 'components/m-browse-actions/m-browse-actions.html'

  bindings:
    path: '<'
    type: '<'

  controller: (playlist, library) ->
    ctrl = this

    stopEvent = (e) ->
      e.preventDefault()
      e.stopPropagation()

    hide = -> ctrl.open = false

    ctrl.showActions = (e) -> stopEvent(e)

    ctrl.add = (e) ->
      playlist.add(ctrl.path, ctrl.type)
      hide()
      stopEvent(e)

    ctrl.addPlay = (e) ->
      playlist.addPlay(ctrl.path, ctrl.type)
      hide()
      stopEvent(e)

    ctrl.addReplacePlay = (e) ->
      playlist.addReplacePlay(ctrl.path, ctrl.type)
      hide()
      stopEvent(e)

    ctrl.update = (e) ->
      library.update(ctrl.path)
      hide()
      stopEvent(e)

    return ctrl
