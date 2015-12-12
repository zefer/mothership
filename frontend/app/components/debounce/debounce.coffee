mod = angular.module('mothership.debounce', [])

mod.service 'debounce', ($timeout, $q) ->
  (fn, wait) ->
    timeout = null
    deferred = $q.defer()

    return ->
      args = arguments
      later = ->
        timeout = null
        result = fn.apply(null, args)
        deferred.resolve(result)
        deferred = $q.defer()

      if timeout
        $timeout.cancel(timeout)
      timeout = $timeout(later, wait, false)

      return deferred.promise
