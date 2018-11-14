mod = angular.module('mothership.airbrake', [])

mod.service 'airbrake', () ->

  # TODO: make these creds configurable.
  airbrake = new airbrakeJs.Client({
    projectId: 203617,
    projectKey: '29aecf58f9f7498977de65220022421b',
    TDigest: tdigest.TDigest,
  })

  airbrake.addFilter (notice) ->
    # TODO: make this configurable and change for each deployment.
    notice.context.environment = 'home'
    notice

  airbrake
