'use strict';

/**
 * @ngdoc function
 * @name virtPanel.controller:MainController
 * @description
 * # MainController
 * Controller of the virtPanel
 */
<<<<<<< Updated upstream
angular.module('virtPanel')
  .controller('HostController', ['$scope', 'Host', function ($scope, Host) {
=======
angular.module('VirtPanel')
  .controller('HostCtrl', ['$scope', 'Host', function ($scope, Host) {
>>>>>>> Stashed changes
    $scope.host = Host.get();
    $scope.containers = [];
  }]);
