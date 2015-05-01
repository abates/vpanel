'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:ContainerCtrl
 * @description
 * # ContainerCtrl
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('ContainerCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
