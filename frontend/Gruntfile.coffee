module.exports = (grunt) ->

  # configuration
  grunt.initConfig

    config:
      build_dir: 'build'
      dist_dir: 'dist'

    copy:
      main:
        files: [
          expand: true
          cwd: 'app/',
          dest: '<%= config.dist_dir %>/'
          src: [
            'index.html'
            'components/**/*.html'
            'assets/*'
          ]
        ]
      vendor:
        files: [
          expand: true
          flatten: true
          cwd: 'bower_components/',
          dest: '<%= config.dist_dir %>/vendor/'
          src: [
            'angular/angular.min.js'
            'angular/angular.min.js.map'
            'angular-ui-router/release/angular-ui-router.min.js'
            'jquery/dist/jquery.min.js'
            'jquery/dist/jquery.min.map'
            'angular-bootstrap/ui-bootstrap-tpls.min.js'
            'font-awesome/css/font-awesome.min.css'
          ]
        ]
      fonts:
        files: [
          expand: true
          cwd: 'bower_components/font-awesome',
          dest: '<%= config.dist_dir %>/'
          src: [
            'fonts/*'
          ]
        ]

    # grunt less
    less:
      compile:
        src: 'app/app.less'
        dest: '<%= config.dist_dir %>/app.css'
        options:
          plugins: [require('less-plugin-glob')]
          style: 'expanded'

    # grunt coffee
    coffee:
      compile:
        expand: true
        cwd: 'app'
        src: [
          'app.coffee'
          '**/*.coffee'
        ]
        dest: '<%= config.build_dir %>/js'
        ext: '.js'
        options:
          bare: true
          preserve_dirs: true

    # combine js
    concat:
      js:
        options:
          separator: ';'
        src: [
          '<%= config.build_dir %>/js/app.js'
          '<%= config.build_dir %>/js/**/*.js'
        ]
        dest: '<%= config.dist_dir %>/app.js'

    shell:
      toBinData:
        command: 'go-bindata -debug -o ../frontend.go -prefix "dist/" <%= config.dist_dir %>/...'

    # grunt watch (or simply grunt)
    watch:
      html:
        files: ['**/*.html']
        tasks: ['copy:main', 'shell:toBinData']
      less:
        files: '**/*.less'
        tasks: ['less', 'shell:toBinData']
      coffee:
        files: '<%= coffee.compile.src %>'
        tasks: ['coffee', 'concat:js', 'shell:toBinData']
      options:
        livereload: true

  # load plugins
  grunt.loadNpmTasks 'grunt-contrib-less'
  grunt.loadNpmTasks 'grunt-contrib-coffee'
  grunt.loadNpmTasks 'grunt-contrib-watch'
  grunt.loadNpmTasks 'grunt-contrib-copy'
  grunt.loadNpmTasks 'grunt-contrib-concat'
  grunt.loadNpmTasks 'grunt-shell'

  # tasks
  grunt.registerTask 'build',   ['copy', 'less', 'coffee', 'concat', 'shell']
  grunt.registerTask 'default', ['build', 'watch']
