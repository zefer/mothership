mod = angular.module('mothership', [
  'ng'
  'ui.router'
  'ui.bootstrap'
  'mothership.airbrake'
  'mothership.mHeader'
  'mothership.mPlaying'
  'mothership.mPlaylist'
  'mothership.mBrowse'
  'mothership.mKeyboard'
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
    url: '/browse?page&sort&direction&filter'
    parent: 'layout'
    template: '<m-browse></m-browse>'

  # This needs to be the last state, because of the wildcard url match.
  $stateProvider.state
    name: 'browse.uri'
    url: '/{uri:nonURIEncoded}?page&sort&direction&filter'
    parent: 'layout'
    template: '<m-browse></m-browse>'

mod.factory '$exceptionHandler', ($log, airbrake) ->
  (exception, cause) ->
    $log.error(exception)
    airbrake.notify({error: exception, params: {angular_cause: cause}})

mod.run ($rootScope, $state, airbrake, mKeyboard) ->
  $rootScope.$state = $state

  $rootScope.$on '$stateChangeStart', (
    _event, toState, _toParams, _fromState, _fromParams
  ) ->
    toState.data = {} if !toState.data
    toState.data.start = new Date()

  $rootScope.$on '$stateChangeSuccess', (event, toState, toParams, fromState, fromParams) ->
    airbrake.incRequest(
      method: "GET",
      route: toState.url,
      statusCode: 200,
      start: toState.data.start,
      end: new Date()
    )

  $rootScope.$on '$stateChangeError', (event, toState, toParams, fromState, fromParams) ->
    airbrake.incRequest(
      method: "GET",
      route: toState.url,
      statusCode: 500,
      start: toState.data.start,
      end: new Date()
    )

  $rootScope.$on '$stateNotFound', (event, toState, toParams, fromState, fromParams) ->
    airbrake.incRequest(
      method: "GET",
      route: toState.url,
      statusCode: 404,
      start: toState.data.start,
      end: new Date()
    )
