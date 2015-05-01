'use strict';

/**
 * @ngdoc service
 * @name virtPanel.Host
 * @description
 * # Host
 * Factory in the virtPanel.
 */
angular.module('virtPanel')
  .factory('Host', [ '$resource', function ($resource) {
    return $resource('host.json');
  }]);
