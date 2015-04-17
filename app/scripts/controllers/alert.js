'use strict';

/**
 * @ngdoc function
 * @name VirtPanel.controller:AlertCtrl
 * @description
 * # AlertCtrl
 * Controller of the VirtPanel
 */
angular.module('VirtPanel')
  .controller('AlertCtrl', ['$scope', 'alerts', function ($scope, alerts) {
    $scope.alerts = alerts;
  }]);
