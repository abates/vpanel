'use strict';

/**
 * @ngdoc directive
 * @name virtPanel.directive:vpFormFeedback
 * @description
 * # vpFormFeedback
 */
angular.module('virtPanel')
  .directive('vpFormFeedback', function () {
    return {
      restrict: 'A',
      compile: function(element, attrs) {
        var formName = attrs.name
        angular.forEach(element.find('input'), function(element) {
          /* It doesn't make sense to validate checkboxes */
          if (element.type === 'checkbox') {
            return;
          }
          angular.element(element).attr('vp-input-feedback', formName);
        });
        angular.element(element[0].getElementsByClassName('form-group')).attr('vp-form-group-class', formName);
      }
    };
  })
  .directive('vpFormGroupClass', function() {
    return {
      restrict: 'A',
      link: function(scope, element, attrs) {
        var formName = attrs.vpFormGroupClass;
        var inputElement = element.find('input')[0]
        if (inputElement === undefined || inputElement === null) {
          return
        }

        var inputName = inputElement.name;
        scope.$watch(formName + '.' + inputName + '.$modelValue', function(value) {
          if (scope[formName] === undefined || 
              scope[formName] === null || 
              scope[formName][inputName] === undefined || 
              scope[formName][inputName] === null) {
            return;
          }

          if (scope[formName][inputName].$dirty) {
            if (!element.hasClass('has-feedback')) {
              element.addClass('has-feedback');
            }

            if (scope[formName][inputName].$valid && !element.hasClass('has-success')) {
              element.removeClass('has-error');
              element.addClass('has-success');
            } else if (scope[formName][inputName].$invalid && !element.hasClass('has-error')) {
              element.removeClass('has-success');
              element.addClass('has-error');
            }
          }
        })
      }
    }
  })
  .directive('vpInputFeedback', function() {
    return {
      restrict: 'A',
      compile: function(element, attrs) {
        var formName = attrs.vpInputFeedback
        var inputName = attrs.name;
        var model = formName + '.' + inputName;

        element.after('<span ng-if="' + model + '.$valid" class="sr-only">' + inputName + ' field is valid</span>');
        element.after('<span ng-if="' + model + '.$dirty && ' + model + '.$invalid" class="sr-only">' + name + ' field is invalid</span>');
        element.after('<span ng-if="' + model + '.$valid" class="glyphicon glyphicon-ok form-control-feedback" aria-hidden="true"></span>');
        element.after('<span ng-if="' + model + '.$dirty && ' + model + '.$invalid" class="glyphicon glyphicon-remove form-control-feedback" aria-hidden="true"></span>');
      }
    }
  });
