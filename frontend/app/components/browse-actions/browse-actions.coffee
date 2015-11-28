mod = angular.module('mothership')

mod.directive 'browseActions', ->
  restrict: 'E'
  templateUrl: 'components/browse-actions/browse-actions.html'

  scope:
    path: '='
    type: '='

  controllerAs: 'ctrl'

  controller: ($scope, playlist, library) ->
    vm = this

    stopEvent = (e) ->
      e.preventDefault()
      e.stopPropagation()

    hide = -> vm.open = false

    vm.showActions = (e) -> stopEvent(e)

    vm.add = (e) ->
      playlist.add($scope.path, $scope.type)
      hide()
      stopEvent(e)

    vm.addPlay = (e) ->
      playlist.addPlay($scope.path, $scope.type)
      hide()
      stopEvent(e)

    vm.addReplacePlay = (e) ->
      playlist.addReplacePlay($scope.path, $scope.type)
      hide()
      stopEvent(e)

    vm.update = (e) ->
      library.update($scope.path)
      hide()
      stopEvent(e)
