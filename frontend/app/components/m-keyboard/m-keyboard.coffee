# Keyboard navigation & controls.

mod = angular.module('mothership.mKeyboard', [
  'mothership.player'
  'mothership.keyboard'
])

mod.service 'mKeyboard', ($rootScope, $state, player, keyboard) ->

  unregisterKeydown = keyboard.onKeydown (event) ->
    # Don't trigger on CMD+/CTRL+ key combinations e.g. CMD+R to reload.
    return if event.metaKey or event.ctrlKey

    switch event.keyCode

      # PLAYER CONTROLS.

      # 32=space, 80=p
      when 32, 80
        if player.state == 'play' then player.pause() else player.play()
      # 39=right-arrow
      when 39
        player.next()
      # 37=left-arrow
      when 37
        player.previous()
      # 82=r
      when 82
        player.random()

      # NAVIGATION.

      # 49=1
      when 49
        $state.go('playing')
      # 50=2
      when 50
        $state.go('browse')
      # 51=3
      when 51
        $state.go('playlist')

      # MISC.

      # 191=?
      when 191
        # Shift+?: Show a list of available shortcuts.
        return unless event.shiftKey

        msg = "Player controls"
        msg += "\n---------------"
        msg += "\nspace  play/pause"
        msg += "\np     play/pause"
        msg += "\n→    next"
        msg += "\n←    previous"
        msg += "\nr     random on/off"
        msg += "\n\n"
        msg += "Navigation"
        msg += "\n-----------"
        msg += "\n1     Playing"
        msg += "\n2     Browse"
        msg += "\n3     Playlist"
        alert(msg)

      else return
    event.preventDefault()

  $rootScope.$on '$destroy', ->
    unregisterKeydown()
