describe 'remember', ->
  beforeEach module('mothership')

  remember = null

  beforeEach inject (_remember_) ->
    remember = _remember_
    localStorage.clear()

  describe 'set', ->
    it 'writes the key/value to localStorage', ->
      remember.set('animal', 'gorilla')
      expect(localStorage.getItem('animal')).to.eq 'gorilla'

  describe 'get', ->
    context 'when the value exists in localStorage', ->
      beforeEach -> localStorage.setItem('food', 'banana')

      it 'reads and returns the value', ->
        expect(remember.get('food', 'orange')).to.eq 'banana'

    context 'when the value exists in localStorage', ->
      it 'reads and returns the value', ->
        expect(remember.get('food', 'orange')).to.eq 'orange'
