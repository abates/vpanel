'use strict';

/**
 * @ngdoc function
 * @name VirtPanel.controller:NetworkingCtrl
 * @description
 * # NetworkingCtrl
 * Controller of the VirtPanel
 */
angular.module('VirtPanel')
  .controller('NetworkingCtrl', ['$scope', 'Bridge', function ($scope, Bridge) {
    $scope.network = {
      utilization: 20
    };

    $scope.bridges = Bridge.query();
    /*$scope.bridges = [{
      'name': 'bridge1',
      'ports': [{
        'vlan': 1,  
        'container': 'container1'
      },{
        'vlan': 3,  
        'container': 'container2'
      },{
        'vlan': 2,  
        'container': 'container3'
      }
    ]},{
      'name': 'bridge2',
      'ports': [{
        'vlan': 7,  
        'container': 'container4'
      },{
        'vlan': 7,  
        'container': 'container5'
      },{
        'vlan': 6,  
        'container': 'container6'
      }
    ]}];*/
  }]);