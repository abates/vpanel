'use strict';

angular.module('VirtPanel', ['ngAnimate', 'ngResource', 'ngRoute'])

  .constant('version', 'v0.1.0')

  .config(function($locationProvider, $routeProvider) {

    $locationProvider.html5Mode(false);

    $routeProvider
      .when('/', {
        templateUrl: 'views/host.html',
        controller: 'HostCtrl'
      })
      .when('/container/:containerId', {
        templateUrl: 'views/container.html',
        controller: 'ContainerCtrl'
      })
      .when('/networking', {
        templateUrl: 'views/networking.html',
        controller: 'NetworkingCtrl'
      })
      .when('/tokens', {
        templateUrl: 'views/tokens.html',
        controller: 'TokensCtrl'
      })
      .when('/audit', {
        templateUrl: 'views/about.html',
        controller: 'AuditCtrl'
      })
      .when('/users', {
        templateUrl: 'views/users.html',
        controller: 'UsersCtrl'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });
