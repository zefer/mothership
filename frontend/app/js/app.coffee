mod = angular.module 'player', [
  'ui.router'
  'ui.bootstrap'
]

mod.config ($stateProvider, $urlRouterProvider) ->

  $urlRouterProvider.otherwise('/playing')

  $stateProvider.state
    name: 'main'
    abstract: true
    views:
      'main':
        controller: 'PlayerCtrl as playerCtrl'
        templateUrl: 'partials/app.html'

  $stateProvider.state
    name: 'playing'
    url: '/playing'
    parent: 'main'
    templateUrl: 'partials/playing.html'

  $stateProvider.state
    name: 'playlist',
    url: '/playlist'
    parent: 'main'
    controller: 'PlaylistCtrl as playlistCtrl'
    templateUrl: 'partials/playlist.html'

  $stateProvider.state
    name: 'browse',
    url: '/browse'
    parent: 'main'
    controller: 'BrowseCtrl as browseCtrl'
    templateUrl: 'partials/browse.html'

  # last state, because of the wildcard url match
  $stateProvider.state
    name: 'browse.uri',
    url: '/{uri:.*}?page&sort&direction'
    parent: 'main'
    controller: 'BrowseCtrl as browseCtrl'
    templateUrl: 'partials/browse.html'

mod.run ($rootScope, $state) ->
  $rootScope.$state = $state
