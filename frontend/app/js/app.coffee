(->
  angular.module('player', ['ui.router', 'ui.bootstrap']).config (
    $stateProvider, $urlRouterProvider) ->
    'use strict'

    $urlRouterProvider.otherwise("/playing")

    $stateProvider.state("playing"
      url: "/playing"
      templateUrl: "partials/playing.html"

    ).state("browse",
      url: "/browse"
      templateUrl: "partials/browse.html"

    ).state("browse.uri",
      url: "/{uri:.*}"
      templateUrl: "partials/browse.html"

    ).state("playlist",
      url: "/playlist"
      templateUrl: "partials/playlist.html"
    )

)().run ($rootScope, $state) ->
  $rootScope.$state = $state
