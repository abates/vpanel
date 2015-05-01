'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('MainCtrl', ['$scope', '$location', function ($scope, $location) {
    $scope.$path = function() {
      return $location.path()
    }
  }]);
