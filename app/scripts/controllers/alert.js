'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:AlertController
 * @description
 * # AlertController
 * Controller of the virtPanel
 */
angular.module('virtPanel')
  .controller('AlertController', ['$scope', 'alerts', function ($scope, alerts) {
    $scope.alerts = alerts;
  }]);
