'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('HostCtrl', ['$scope', 'Host', 'Container', function ($scope, Host, Container) {
    //$scope.host = Host.get();
    $scope.host = {};
    //$scope.containers = Container.query();
    $scope.containers = [];
  }]);
