'use strict';

/**
 * @ngdoc service
 * @name VirtPanel.bridge
 * @description
 * # bridge
 * Factory in the VirtPanel.
 */
angular.module('VirtPanel')
  .factory('Bridge', ['$resource', function ($resource) {
    var Bridge = $resource('bridge_:id.json', {}, {
      query: {method: 'GET', url: 'bridges.json', isArray: true}
    });

    return Bridge;
  }]);
