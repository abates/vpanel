'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainController
 * @description
 * # MainController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('MainController', ['$scope', '$location', function ($scope, $location) {
    $scope.$path = function() {
      return $location.path()
    }
  }]);
