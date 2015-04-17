'use strict';

var gulp = require('gulp');
var config = require('ng-factory').use(gulp);
var portfinder = require('portfinder');
var symlink = require('gulp-symlink');

portfinder.basePort = 9000;

//
// Aliases
gulp.task('serve', gulp.series('ng:serve'));
gulp.task('build', gulp.series('ng:build'));

var serverPort;

gulp.task('go:symlink:bower', function() {
  return gulp.src('app/bower_components').pipe(symlink('.tmp/bower_components'))
});

gulp.task('go:symlink:images', function() {
  return gulp.src('app/images').pipe(symlink('.tmp/images'))
});

gulp.task('go:src/serve', function(cb) {
  var http = require('http');

  function checkServer(port, cb) {
    setTimeout(function () {
      http.request({
        method: 'HEAD',
        hostname: 'localhost',
        port: port
      }, function (res) {
        return cb();
      }).on('error', function (err) {
        checkServer(port, cb);
      }).end();
    }, 50);
  }

  portfinder.getPorts(2, {}, function(err, ports) {
    if (err) {
      console.log("go:src/serve could not allocate a free port: ", err)
      return
    }

    serverPort = ports[0];
    var appPort = ports[1];

    var spawn = require('child_process').spawn
    
    // rebuild bin assets
    var bindata = spawn('go-bindata', ['-debug', '-prefix', __dirname + '/.tmp/', '.tmp/...'], { stdio: 'inherit' })

    bindata.on('exit', function(code) {
      // Start gin server
      spawn('gin', ['--port', serverPort, '--appPort', appPort, '-i', 'run'], { stdio: 'inherit' })
      checkServer(appPort, cb);
    })
  });
});

gulp.task('go:proxy', function(cb) {
  var browserSync = require('browser-sync');
  portfinder.getPort(function(err, port) {
    if (err) {
      console.log("go:proxy could not allocate a free port: ", err)
      return
    }

    browserSync({
      ui: false,
      port: port,
      notify: false,
      open: true,
      logPrefix: function () {
        return this.compile('[{gray:' + new Date().toLocaleTimeString() + '}] ');
      },
      proxy: 'localhost:' + serverPort
    }, cb);
  });
});

gulp.set('go:symlink', gulp.series('go:symlink:bower', 'go:symlink:images'));
gulp.set('go:serve', gulp.series('ng:beforeServe', 'ng:src/clean', 'ng:src/views', 'go:symlink', ['go:src/serve', 'go:proxy', 'ng:src/watch'], 'ng:afterServe'));

//
// Hooks example

// var path = require('path');
// var src = config.src;
// gulp.task('ng:afterBuild', function() {
//   gulp.src(['bower_components/font-awesome/fonts/*.woff'], {cwd: src.cwd})
//     .pipe(gulp.dest(path.join(src.dest, 'fonts')));
// });
