mod = angular.module 'mothership', [
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
    template: '<m-player-status></m-player-status>'

  $stateProvider.state
    name: 'playlist'
    url: '/playlist'
    parent: 'main'
    template: '<m-playlist></m-playlist>'

  $stateProvider.state
    name: 'browse',
    url: '/browse?page&sort&direction'
    parent: 'main'
    template: '<m-browse></m-browse>'

  # last state, because of the wildcard url match
  $stateProvider.state
    name: 'browse.uri'
    url: '/{uri:.*}?page&sort&direction'
    parent: 'main'
    template: '<m-browse></m-browse>'

mod.run ($rootScope, $state) ->
  $rootScope.$state = $state
