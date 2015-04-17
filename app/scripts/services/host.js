'use strict';

/**
 * @ngdoc service
 * @name VirtPanel.Host
 * @description
 * # Host
 * Factory in the VirtPanel.
 */
angular.module('VirtPanel')
  .factory('Host', [ '$resource', function ($resource) {
    return $resource('host.json');
  }]);
