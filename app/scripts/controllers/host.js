'use strict';

/**
 * @ngdoc function
 * @name VirtPanel.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the VirtPanel
 */
angular.module('VirtPanel')
  .controller('HostCtrl', ['$scope', 'Host', 'Container', function ($scope, Host, Container) {
    //$scope.host = Host.get();
    $scope.host = {};
    //$scope.containers = Container.query();
    $scope.containers = [];
  }]);
