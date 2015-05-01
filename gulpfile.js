'use strict';

var gulp = require('gulp');
var config = require('ng-factory').use(gulp);
var url = require('url');

//
// Aliases
gulp.task('serve', gulp.series('ng:serve'));
gulp.task('build', gulp.series('ng:build'));

var proxy = require('proxy-middleware');
var proxyOptions = url.parse('http://192.168.56.101:3000/api');
proxyOptions.route = '/api';
config.middleware = [proxy(proxyOptions)];

