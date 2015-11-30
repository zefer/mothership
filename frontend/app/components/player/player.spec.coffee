describe 'player', ->
  beforeEach module('mothership.player')

  player = null
  mpd = null

  beforeEach inject (_player_, _mpd_) ->
    player = _player_
    mpd = _mpd_

  describe 'play()', ->
    it 'delegates to the mpd service', ->
      sinon.spy(mpd, 'play')
      player.play()
      expect(mpd.play).to.have.been.called

  describe 'pause()', ->
    it 'delegates to the mpd service', ->
      sinon.spy(mpd, 'pause')
      player.pause()
      expect(mpd.pause).to.have.been.called

  describe 'previous()', ->
    it 'delegates to the mpd service', ->
      sinon.spy(mpd, 'previous')
      player.previous()
      expect(mpd.previous).to.have.been.called

  describe 'next()', ->
    it 'delegates to the mpd service', ->
      sinon.spy(mpd, 'next')
      player.next()
      expect(mpd.next).to.have.been.called

  describe 'random()', ->
    context 'when random is enabled', ->
      beforeEach -> player.randomOn = true

      it 'calls randomOff() on the mpd service', ->
        sinon.spy(mpd, 'randomOff')
        player.random()
        expect(mpd.randomOff).to.have.been.called

    context 'when random is disabled', ->
      beforeEach -> player.randomOn = false

      it 'calls randomOn() on the mpd service', ->
        sinon.spy(mpd, 'randomOn')
        player.random()
        expect(mpd.randomOn).to.have.been.called

