'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('AboutCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
