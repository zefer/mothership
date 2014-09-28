# Calculates the elapsed time for the current track & updates the scope directly
# to allow the UI to update.

mod = angular.module("player")

mod.factory "progress", ["$interval", ($interval) ->
  'use strict'
  that     = this
  timer    = null
  scope    = null
  duration = 1000

  end     = null # time
  total   = null # seconds
  elapsed = null # seconds

  that.updateScope = ->
    scope.percentage = Math.floor((elapsed/total)*100)
    scope.elapsed    = that.secondsToString(elapsed)
    scope.total      = that.secondsToString(total)

  that.startTimer = ->
    timer ?= $interval ->
      cur = (new Date()).getTime()
      remaining = Math.floor((end-cur)/1000)
      elapsed = total - remaining
      that.updateScope()
    , duration

  that.stopTimer = ->
    timer && $interval.cancel(timer)
    timer = null

  that.secondsToString = (secs) ->
    mins = Math.floor(secs / 60)
    secs -= mins * 60
    mm = if mins < 10 then '0' + mins else mins
    ss = if secs < 10 then '0' + secs else secs
    mm + ':' + ss

  api =
    update: (data, scp) ->
      scope = scp
      total = parseInt(data.Time)
      end = (new Date()).getTime() + total*1000 - parseInt(data.elapsed)*1000
      if isNaN(total) || isNaN(end) || data.state != "play"
        that.stopTimer()
      else
        that.startTimer()

  api
]
