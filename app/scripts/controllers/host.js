'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainController
 * @description
 * # MainController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('HostController', ['$scope', 'Host', 'Container', function ($scope, Host, Container) {
    //$scope.host = Host.get();
    $scope.host = {};
    //$scope.containers = Container.query();
    $scope.containers = [];
  }]);
