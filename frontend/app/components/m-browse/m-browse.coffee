mod = angular.module('mothership.mBrowse', [
  'ui.router'

  'mothership.library'
  'mothership.remember'
  'mothership.mSortBy'
  'mothership.mPagination'
  'mothership.mBrowseActions'
  'mothership.mSearch'
])

mod.component 'mBrowse',
  bindings: {}
  templateUrl: 'components/m-browse/m-browse.html'

  controller: ($scope, $state, library, remember) ->
    ctrl = this

    MAX_PER_PAGE = 200

    list = ->
      params = $state.params
      params.filter ?= ''
      params.uri ?= '/'
      params.uri = '/' if params.uri == ''
      params.page ?= 1
      params.sort ?= remember.get('sort', 'date')
      params.direction ?= remember.get('direction', 'desc')
      remember.set('sort', params.sort)
      remember.set('direction', params.direction)
      library.ls(
        params.uri, params.sort, params.direction, params.filter
      ).then (items) ->
        paginate(items, parseInt(params.page))
        ctrl.breadcrumbs = breadcrumbs(params.uri)

    do events = ->
      $scope.$on 'search:filter', ->
        list()

    do init = ->
      list()

    breadcrumbs = (uri) ->
      parts = uri.split('/')
      crumbs = ({
        label: part,
        path: parts[0..i].join('/')
      } for part, i in parts when part != '')
      crumbs.unshift { label: 'home', path: '' }
      crumbs

    paginate = (items, page) ->
      if !items? or items.length < 1
        ctrl.pages = []
        ctrl.items = []
        return

      pages = Math.ceil(items.length / MAX_PER_PAGE)
      page = pages if page > pages
      page = 1 if page < 1
      pos = (page - 1) * MAX_PER_PAGE
      ctrl.items = items[pos..pos+MAX_PER_PAGE]
      ctrl.pages = (i for i in [1..pages])
      ctrl.page = page

    return ctrl
