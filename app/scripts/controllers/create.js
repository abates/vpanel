'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:CreateController
 * @description
 * # CreateController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('CreateController', ['$scope', 'Container', function($scope, Container) {
    $scope.loading = false;
    $scope.error = false;
    $scope.errorMessage = "";

    $scope.templates = Container.templates();
    $scope.container = {
      name: '',
      hostname: '',
      template: '',
      autostart: false
    };

    $scope.submit = function(form) {
      if (form.$invalid) {
        console.log("Cannot submit an invalid form!");
        return;
      }

      $scope.loading = true;
      Container.save($scope.container, function() {
        $scope.loading = false;
        $scope.$hide();
      }, function(resp) {
        $scope.loading = false;
        $scope.error = true;
        console.log(resp);
        if (angular.isString(resp.data)) {
          $scope.errorMessage = resp.data;
        } else {
          $scope.errorMessage = resp.data.message;
        }
      });
    };
  }]);
