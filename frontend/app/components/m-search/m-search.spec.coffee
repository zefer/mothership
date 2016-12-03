markup = '<m-search></m-search>'
$rootScope = null
$state = null
$timeout = null
ctrl = null

describe 'mSearch', ->
  beforeEach module('mothership.mSearch')
  beforeEach module('mothership.templates')

  beforeEach inject (
    _$rootScope_, _$state_, _$timeout_, $compile
  ) ->
    $rootScope = _$rootScope_
    scope = $rootScope.$new()
    $state = _$state_
    $timeout = _$timeout_

    elem = $compile(markup)(scope)
    scope.$digest()
    ctrl = elem.controller('mSearch')

    sinon.stub($state, 'go')
    sinon.spy($rootScope, '$broadcast')

  describe 'search', ->
    beforeEach ->
      ctrl.filter = 'elvis'
      ctrl.search()
      $timeout.flush()

    it 'updates "filter" in state params with the new filter text', ->
      expect($state.params.filter).to.eq('elvis')

    it 'navigates to the new state', ->
      expect($state.go).to.have.been.calledWithExactly(
        '.', {filter: 'elvis'}, {notify: false}
      )

    it 'broadcasts the "search:filter" event', ->
      expect($rootScope.$broadcast).to.have.been.calledWith('search:filter')

  describe 'clear', ->
    beforeEach ->
      ctrl.filter = 'elvis'
      $state.params.filter = ctrl.filter
      ctrl.clear()
      $timeout.flush()

    it 'clears the "filter" input', ->
      expect(ctrl.filter).to.eq('')

    it 'clears "filter" in state params', ->
      expect($state.params.filter).to.eq('')

    it 'navigates to the state with an empty filter', ->
      expect($state.go).to.have.been.calledWithExactly(
        '.', {filter: ''}, {notify: false}
      )

    it 'broadcasts the "search:filter" event', ->
      expect($rootScope.$broadcast).to.have.been.calledWith('search:filter')
