'use strict';

/**
 * @ngdoc service
 * @name virtPanel.bridge
 * @description
 * # bridge
 * Factory in the virtPanel.
 */
angular.module('virtPanel')
  .factory('Bridge', ['$resource', function ($resource) {
    var Bridge = $resource('bridge_:id.json', {}, {
      query: {method: 'GET', url: 'bridges.json', isArray: true}
    });

    return Bridge;
  }]);
