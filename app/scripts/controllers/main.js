'use strict';

/**
 * @ngdoc function
 * @name VirtPanel.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the VirtPanel
 */
angular.module('VirtPanel')
  .controller('MainCtrl', ['$scope', '$location', function ($scope, $location) {
    $scope.$path = function() {
      return $location.path()
    }
  }]);
