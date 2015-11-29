mod = angular.module 'mothership'

mod.directive 'mNavigation', ->
  restrict: 'A'
  scope: {}
  templateUrl: 'components/m-navigation/m-navigation.html'

  controller: ($scope, $state, playlist) ->
    vm = this

    $scope.playlist = playlist

    $scope.active = (tabName) ->
      active = if $state.includes('playing')
        'playing'
      else if $state.includes('playlist')
        'playlist'
      else
        'browse'
      tabName == active
