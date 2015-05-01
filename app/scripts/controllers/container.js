'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:ContainerController
 * @description
 * # ContainerController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('ContainerController', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
