module.exports = (grunt) ->

  # configuration
  grunt.initConfig

    config:
      build_dir: 'build'
      dist_dir: 'dist'

    connect:
      server:
        options:
          port: 4000,
          base: './dist'

    copy:
      main:
        files: [
          expand: true
          cwd: 'app/',
          dest: '<%= config.dist_dir %>/'
          src: [
            'index.html',
            'partials/*'
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
            'jquery/dist/jquery.min.js'
            'jquery/dist/jquery.min.map'
            'bootstrap/dist/css/bootstrap.min.css'
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
        options:
          style: 'expanded'
        files: [
          expand: true
          cwd: 'app/css'
          src: ['**/*.less']
          dest: '<%= config.build_dir %>/css'
          ext: '.css'
        ]

    # grunt coffee
    coffee:
      compile:
        expand: true
        cwd: 'app/js'
        src: ['**/*.coffee']
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
      css:
        src: ['<%= config.build_dir %>/css/**/*.css']
        dest: '<%= config.dist_dir %>/style.css'

    # grunt watch (or simply grunt)
    watch:
      html:
        files: ['**/*.html']
        tasks: ['copy:main']
      less:
        files: '<%= less.compile.files[0].src %>'
        tasks: ['less', 'concat:css']
      coffee:
        files: '<%= coffee.compile.src %>'
        tasks: ['coffee', 'concat:js']
      options:
        livereload: true

  # load plugins
  grunt.loadNpmTasks 'grunt-contrib-less'
  grunt.loadNpmTasks 'grunt-contrib-coffee'
  grunt.loadNpmTasks 'grunt-contrib-watch'
  grunt.loadNpmTasks 'grunt-contrib-copy'
  grunt.loadNpmTasks 'grunt-contrib-connect'
  grunt.loadNpmTasks 'grunt-contrib-concat'

  # tasks
  grunt.registerTask 'default', ['connect', 'copy', 'less', 'coffee', 'concat', 'watch']
