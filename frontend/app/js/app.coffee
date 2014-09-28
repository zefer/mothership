(->
  angular.module('player', ['ui.router', 'ui.bootstrap']).config (
    $stateProvider, $urlRouterProvider) ->
    'use strict'

    $urlRouterProvider.otherwise("/playing")

    $stateProvider.state("main"
      abstract: true
      views:
        "main":
          controller: "PlayerCtrl as playerCtrl"
          templateUrl: "partials/app.html"

    ).state("playing"
      url: "/playing"
      parent: "main"
      templateUrl: "partials/playing.html"

    ).state("playlist",
      url: "/playlist"
      parent: "main"
      controller: "PlaylistCtrl as playlistCtrl"
      templateUrl: "partials/playlist.html"

    ).state("browse",
      url: "/browse"
      parent: "main"
      controller: "BrowseCtrl as browseCtrl"
      templateUrl: "partials/browse.html"

    # last state, because of the wildcard url match
    ).state("browse.uri",
      url: "/{uri:.*}"
      parent: "main"
      controller: "BrowseCtrl as browseCtrl"
      templateUrl: "partials/browse.html"
    )

)().run ($rootScope, $state) ->
  $rootScope.$state = $state
