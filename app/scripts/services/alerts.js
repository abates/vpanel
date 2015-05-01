'use strict';

/**
 * @ngdoc service
 * @name virtPanel.alert
 * @description
 * # alert
 * Factory in the virtPanel.
 */
angular.module('virtPanel')
  .factory('alerts', ['$timeout', function ($timeout) {
    var autoDismiss = null;
    var alerts = [];

    function setupAutoDismiss() {
      if (autoDismiss === null && alerts.length >= 0) {
        autoDismiss = $timeout(function() {
          alerts.shift();
          autoDismiss = null;
          setupAutoDismiss();
        }, 5000);
      }
    }

    var Alerts = {
      push: function(level, title, message) {
        alerts.push({
          classes: {
            'animated': true,
            'slideInDown': true
          },
          title: title,
          message: message
        });
        alerts[0].classes['alert-' + level] = true;
        setupAutoDismiss();
      },

      empty: function() {
        return alerts.length === 0;
      },

      current: function() {
        return alerts[0];
      },

      dismiss: function() {
        // cancel auto dismiss (if any)
        if (autoDismiss !== null) {
          $timeout.cancel(autoDismiss);
        }

        alerts.shift();

        // setup auto dismiss for any remaining alerts
        setupAutoDismiss();
      },
    };

    angular.forEach(['success', 'info', 'warning', 'danger'], function(level) {
      Alerts[level] = function(title, message) {
        Alerts.push(level, title, message);
      };
    });

    return Alerts;
  }]);
