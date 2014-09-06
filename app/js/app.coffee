(->
  'use strict'

  angular.module('player', [])

  .controller('PlayerController', ->
  )

  .directive('playerControls', ->
    restrict: 'E'
    templateUrl: 'partials/play-controls.html',
  )

)()
