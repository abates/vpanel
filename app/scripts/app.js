'use strict';

angular.module('virtPanel', ['ngAnimate', 'ngResource', 'ngRoute', 'mgcrea.ngStrap'])

  .constant('version', 'v0.1.0')

  .config(function($locationProvider, $routeProvider) {

    $locationProvider.html5Mode(false);

    $routeProvider
      .when('/', {
        templateUrl: 'views/host.html',
        controller: 'HostController'
      })
      .when('/container/:containerId', {
        templateUrl: 'views/container.html',
        controller: 'ContainerController'
      })
      .when('/networking', {
        templateUrl: 'views/networking.html',
        controller: 'NetworkingController'
      })
      .when('/tokens', {
        templateUrl: 'views/tokens.html',
        controller: 'TokensController'
      })
      .when('/audit', {
        templateUrl: 'views/about.html',
        controller: 'AuditController'
      })
      .when('/users', {
        templateUrl: 'views/users.html',
        controller: 'UsersController'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutController'
      })
      .otherwise({
        redirectTo: '/'
      });
  });
