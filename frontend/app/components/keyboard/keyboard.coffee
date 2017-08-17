mod = angular.module('mothership.keyboard', [])


# Callbacks will not fire when the focus is in these element types, this is to
# allow ordinary keyboard input with form elements.
DISABLE_WHEN_FOCUS = ['search', 'text', 'textarea']

mod.service 'keyboard', ($document, $window, $rootScope) ->
  callbacks = keydown: [], keypress: []
  disabled  = false

  service =
    # Explosing this is a small hack, so we can stub it in tests. An alternative
    # solution is to make this a separate service and inject it as a dependency.
    focusElementType: -> $document[0].activeElement.type

  $document.on 'keypress', (event) -> keypress(event)
  $document.on 'keydown', (event) -> keydown(event)

  keydown = (event) ->
    # Disable when typing in normal input elements such as text or textarea,
    # to allow normal input.
    disabled = (DISABLE_WHEN_FOCUS.indexOf(service.focusElementType()) >= 0)
    # make ng aware of changes outside the normal working loop
    $rootScope.$apply ->
      cb(event) for cb in callbacks.keydown unless disabled

  keypress = (event) ->
    # make ng aware of changes outside the normal working loop
    $rootScope.$apply ->
      cb(event) for cb in callbacks.keypress unless disabled

  onKeypress = (callback) ->
    callbacks.keypress.push callback
    # Return a function to unregister.
    return ->
      index = callbacks.keypress.indexOf(callback)
      callbacks.keypress.splice(index, 1) if index > -1

  onKeydown = (callback) ->
    callbacks.keydown.push callback
    # Return a function to unregister.
    return ->
      index = callbacks.keydown.indexOf(callback)
      callbacks.keydown.splice(index, 1) if index > -1

  service = angular.extend service,
    # These 2 functions are exposed so we can call them from our tests.
    # Otherwise these would remain private.
    keypress: keypress
    keydown:  keydown

    onKeypress: onKeypress
    onKeydown:  onKeydown
