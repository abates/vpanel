'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:AlertCtrl
 * @description
 * # AlertCtrl
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('AlertCtrl', ['$scope', 'alerts', function ($scope, alerts) {
    $scope.alerts = alerts;
  }]);
