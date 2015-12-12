describe 'mBrowse', ->
  beforeEach module('mothership.mBrowse')
  beforeEach module('mothership.templates')

  $q = null
  $scope = null
  $compile = null
  $state = null
  library = null
  elem = null

  compile = (markup) ->
    elem = $compile(angular.element(markup))($scope)
    $scope.$digest()
    return elem

  files = [{
    path: "somewhere/radio",
    type: "directory",
    base: "radio",
    lastModified: "2015-11-28T15:09:26Z"
  }, {
    path: "somewhere/music",
    type: "directory",
    base: "music",
    lastModified: "2015-11-21T12:08:18Z"
  }]

  beforeEach inject (_$q_, $rootScope, _$compile_, $window, _$state_, _library_) ->
    $q = _$q_
    $compile = _$compile_
    $scope = $rootScope.$new()
    $state = _$state_
    library = _library_

  renderDirective = (files) ->
    sinon.stub(library, 'ls').returns($q.when(angular.copy(files)))
    markup = '<m-browse></m-browse>'
    compile(markup)

  it 'renders the expected directives', ->
    html = renderDirective(files).html()
    # TODO: make breadcrumbs a directive.
    expect(html).to.contain('breadcrumb')
    expect(html).to.contain('<m-sort-by ')
    expect(html).to.contain('<m-pagination pages="pages" page="page" ')
    expect(html).to.contain('<m-browse-actions ')

  it 'fetches the library listing for the current state params', ->
    $state.params.uri = '/somewhere'
    $state.params.sort = 'date'
    $state.params.direction = 'desc'
    renderDirective(files)

    expect(library.ls).to.have.been.calledWithExactly(
      '/somewhere', 'date', 'desc'
    )

  describe 'pagination', ->
    beforeEach ->
      $state.params.page = '2'
      files = ({
        path: "path#{i}",
        type: 'directory',
        base: "base#{i}",
        lastModified: "2015-11-21T12:08:18Z"
      } for i in [1..202])
      expect(files.length).to.eq(202)
      elem = renderDirective(files)

    it 'renders list items from the current page', ->
      expect(elem.html()).to.contain(path) for path in [
        'base201', 'base202'
      ]

    it 'does not render list items from other pages', ->
      expect(elem.html()).not.to.contain(path) for path in [
        'base1', 'base200', 'base203'
      ]

    it 'renders top pagination links for all the pages', ->
      pagination = angular.element(elem.find('m-pagination')[0])
      items = pagination.find('li')
      expect(items.length).to.eq(4)
      expect(angular.element(items[0]).text().trim()).to.eq('«')
      expect(angular.element(items[1]).text().trim()).to.eq('1')
      expect(angular.element(items[2]).text().trim()).to.eq('2')
      expect(angular.element(items[3]).text().trim()).to.eq('»')

    it 'renders bottom pagination links for all the pages', ->
      pagination = angular.element(elem.find('m-pagination')[1])
      items = pagination.find('li')
      expect(items.length).to.eq(4)
      expect(angular.element(items[0]).text().trim()).to.eq('«')
      expect(angular.element(items[1]).text().trim()).to.eq('1')
      expect(angular.element(items[2]).text().trim()).to.eq('2')
      expect(angular.element(items[3]).text().trim()).to.eq('»')

  # TODO: make breadcrumbs a directive.
  describe 'breadcrumbs', ->
    breadcrumbs = null

    context 'when browsing root', ->
      beforeEach ->
        $state.params.uri = '/'
        breadcrumbs = renderDirective(files).find('ol').find('li')

      it 'only renders "home"', ->
        expect(breadcrumbs.length).to.eq(1)
        expect(angular.element(breadcrumbs[0]).text()).to.contain('home')

    context 'when browsing a subpath', ->
      beforeEach ->
        $state.params.uri = '/music/banana/gorilla'
        breadcrumbs = renderDirective(files).find('ol').find('li')

      it 'renders the all the path items, including "home"', ->
        expect(breadcrumbs.length).to.eq(4)
        for label, i in ['home', 'music', 'banana', 'gorilla']
          expect(angular.element(breadcrumbs[i]).text()).to.contain(label)
