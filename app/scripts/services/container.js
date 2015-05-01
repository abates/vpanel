'use strict';

/**
 * @ngdoc service
 * @name virtPanel.Container
 * @description
 * # Container
 * Factory in the virtPanel.
 */
angular.module('virtPanel')
  .factory('Container', ['$resource', 'alerts', function ($resource, alerts) {
    var Container = $resource('containers/:id.json', {}, {
      query: {method: 'GET', url: 'containers.json', isArray: true}
    });

    Container.prototype.start = function() {
      if (this.stat !== 'Running') {
        this.stat = 'Running';
      } else {
        alerts.warning('Invalid State', 'Cannot start a container that is already running');
      }
    };

    Container.prototype.stop = function() {
      if (this.stat !== 'Stopped') {
        this.stat = 'Stopped';
      } else {
        alerts.warning('Invalid State', 'Cannot stop a container that is already stopped');
      }
    };

    Container.prototype.freeze = function() {
      if (this.stat !== 'Frozen') {
        this.stat = 'Frozen';
      } else {
        alerts.warning('Invalid State', 'Cannot freeze a container that is already frozen');
      }
    };

    return Container;
  }]);
