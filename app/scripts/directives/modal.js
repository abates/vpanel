'use strict';

/**
 * @ngdoc directive
 * @name virtPanel.directive:modal
 * @description
 * # modal
 */
angular.module('virtPanel')
  .directive('vpModal', function () {
    return {
      templateUrl: '/views/modals/modal.tpl.html',
      restrict: 'E',
      replace: true,
      transclude: true,
      link : function($scope, element, attrs, controller, transclude) {
        $scope.label = attrs['vpLabel'];
        $scope.labelClass = attrs['vpLabelClass'];
        transclude($scope, function(clone){
          // jqlite w00t!
          angular.element(element[0].getElementsByClassName('transclude')[0]).replaceWith(clone);
        });
      }
    };
  });
