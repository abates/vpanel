'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:TokensController
 * @description
 * # TokensController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('TokensController', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
