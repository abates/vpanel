'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainController
 * @description
 * # MainController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('HostController', ['$scope', 'Host', function ($scope, Host) {
    $scope.host = Host.get();
    $scope.containers = [];
  }]);
