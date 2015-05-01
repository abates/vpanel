'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:AboutController
 * @description
 * # AboutController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('AboutController', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
