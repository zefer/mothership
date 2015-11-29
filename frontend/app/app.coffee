mod = angular.module('mothership', [
  'ng'
  'ui.router'
  'ui.bootstrap'
  'mothership.mHeader'
  'mothership.mPlaying'
  'mothership.mPlaylist'
  'mothership.mBrowse'
])

mod.config (
  $stateProvider, $urlRouterProvider, $urlMatcherFactoryProvider
) ->

  $urlRouterProvider.otherwise('/playing')

  valToString = (val) -> if val? then val.toString() else val
  $urlMatcherFactoryProvider.type 'nonURIEncoded',
    encode: valToString
    decode: valToString
    is: -> true

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
    url: '/{uri:nonURIEncoded}?page&sort&direction'
    parent: 'layout'
    template: '<m-browse></m-browse>'

mod.run ($rootScope, $state) ->
  $rootScope.$state = $state
