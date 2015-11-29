mod = angular.module 'mothership', [
  'ui.router'
  'ui.bootstrap'
]

mod.config ($stateProvider, $urlRouterProvider) ->

  $urlRouterProvider.otherwise('/playing')

  $stateProvider.state
    name: 'layout'
    abstract: true
    views:
      'layout':
        template: '<m-header></m-header><main ui-view></main>'

  $stateProvider.state
    name: 'playing'
    url: '/playing'
    parent: 'layout'
    template: '<m-playing></m-playing>'

  $stateProvider.state
    name: 'playlist'
    url: '/playlist'
    parent: 'layout'
    template: '<m-playlist></m-playlist>'

  $stateProvider.state
    name: 'browse',
    url: '/browse?page&sort&direction'
    parent: 'layout'
    template: '<m-browse></m-browse>'

  # This needs to be the last state, because of the wildcard url match.
  $stateProvider.state
    name: 'browse.uri'
    url: '/{uri:.*}?page&sort&direction'
    parent: 'layout'
    template: '<m-browse></m-browse>'

mod.run ($rootScope, $state) ->
  $rootScope.$state = $state
