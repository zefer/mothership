keyboard = null
event    = null
callback = null

describe 'keyboard', ->
  beforeEach module('mothership.keyboard')

  beforeEach inject (_keyboard_) ->
    keyboard = _keyboard_
    event = keyCode: 74
    callback = sinon.spy()

  describe 'when the focus is not on a known type of text input element', ->
    beforeEach -> keyboard.focusElementType = -> "checkbox"

    it 'fires on keypress with the event', ->
      keyboard.onKeypress(callback)
      keyboard.keypress(event)
      expect(callback).to.have.been.calledWith(event)

    it 'fires on keydown with the event', ->
      keyboard.onKeydown(callback)
      keyboard.keydown(event)
      expect(callback).to.have.been.calledWith(event)

  describe 'after unregistering the callback', ->
    it 'no longer fires on keypress', ->
      unregister = keyboard.onKeypress(callback)
      unregister()
      keyboard.keypress(event)
      expect(callback).not.to.have.been.called

    it 'no longer fires on keydown', ->
      unregister = keyboard.onKeydown(callback)
      unregister()
      keyboard.keydown(event)
      expect(callback).not.to.have.been.called

  describe 'when the focus is on a known type of text input element', ->
    types = ['search', 'text', 'textarea']

    it 'does not fire on keypress', ->
      for type in types
        keyboard.focusElementType = -> type
        keyboard.onKeypress(callback)
        # Fire keydown first, this will always fire before keypress & is needed.
        keyboard.keydown(event)
        keyboard.keypress(event)
        expect(callback).not.to.have.been.called

    it 'does not fire on keydown', ->
      for type in types
        keyboard.focusElementType = -> type
        keyboard.onKeydown(callback)
        keyboard.keydown(event)
        expect(callback).not.to.have.been.called
