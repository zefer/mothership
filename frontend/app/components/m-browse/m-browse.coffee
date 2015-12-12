mod = angular.module('mothership.mBrowse', [
  'ui.router'

  'mothership.library'
  'mothership.remember'
  'mothership.mSortBy'
  'mothership.mPagination'
  'mothership.mBrowseActions'
])

mod.directive 'mBrowse', ->
  restrict: 'E'
  scope: {}
  templateUrl: 'components/m-browse/m-browse.html'

  controller: ($scope, $state, library, remember) ->
    vm = this

    MAX_PER_PAGE = 200

    do init = ->
      params = $state.params
      params.uri ?= '/'
      params.uri = '/' if params.uri == ''
      params.page ?= 1
      params.sort ?= remember.get('sort', 'date')
      params.direction ?= remember.get('direction', 'desc')
      remember.set('sort', params.sort)
      remember.set('direction', params.direction)
      library.ls(params.uri, params.sort, params.direction).then (items) ->
        paginate(items, parseInt(params.page))
        $scope.breadcrumbs = breadcrumbs(params.uri)

    breadcrumbs = (uri) ->
      parts = uri.split('/')
      crumbs = ({
        label: part,
        path: parts[0..i].join('/')
      } for part, i in parts when part != '')
      crumbs.unshift { label: 'home', path: '' }
      crumbs

    paginate = (items, page) ->
      pages = Math.ceil(items.length / MAX_PER_PAGE)
      page = pages if page > pages
      page = 1 if page < 1
      pos = (page - 1) * MAX_PER_PAGE
      $scope.items = items[pos..pos+MAX_PER_PAGE]
      $scope.pages = (i for i in [1..pages])
      $scope.page = page
