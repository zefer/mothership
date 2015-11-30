module.exports = function (config) {
  'use strict';
  config.set({

    basePath: '',

    frameworks: ['mocha', 'chai', 'sinon-chai'],

    files: [
      'bower_components/angular/angular.min.js',
      'bower_components/angular-mocks/angular-mocks.js',
      'bower_components/angular-ui-router/release/angular-ui-router.min.js',
      'bower_components/angular-bootstrap/ui-bootstrap.min.js',
      'build/js/**/*.js',
      'build/spec/**/*.js'
    ],

    preprocessors: {
      'build/spec/**/*.js': ['sourcemap']
    },

    reporters: ['spec'],

    port: 9876,
    colors: true,
    autoWatch: false,
    singleRun: false,

    // level of logging
    // possible values:
    // config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN ||
    // config.LOG_INFO || config.LOG_DEBUG
    logLevel: config.LOG_INFO,

    browsers: ['PhantomJS']
  });
};
