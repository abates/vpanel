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
    };

    $scope.inputClass = function(input) {
      if (input === undefined || input === null) {
        return "";
      }

      if (input.$dirty) {
        if (input.$valid) {
          return "has-success has-feedback";
        } else if (input.$invalid) {
          return "has-error has-feedback";
        }
      }
      return "";
    };
  }]);
