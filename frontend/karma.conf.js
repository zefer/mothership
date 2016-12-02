module.exports = function (config) {
  'use strict';
  config.set({

    basePath: '',

    frameworks: ['mocha', 'chai', 'sinon-chai'],

    files: [
      'node_modules/angular/angular.min.js',
      'node_modules/angular-mocks/angular-mocks.js',
      'node_modules/angular-ui-router/release/angular-ui-router.min.js',
      'node_modules/angular-ui-bootstrap/dist/ui-bootstrap.js',
      'build/js/**/*.js',
      'build/spec/**/*.js',
      'dist/**/*.html'
    ],

    preprocessors: {
      'build/spec/**/*.js': ['sourcemap'],
      'dist/**/*.html': ['ng-html2js']
    },

    ngHtml2JsPreprocessor: {
      stripPrefix: 'dist/',
      moduleName: "mothership.templates"
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
