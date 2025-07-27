#ifndef ADVANCED_CALCULATOR_H
#define ADVANCED_CALCULATOR_H

#include "calculator.h"  // Demonstrates dependency
#include <vector>

// An advanced calculator that uses the basic calculator.
class AdvancedCalculator {
   public:
    // Computes the power of a number.
    double power(double base, int exp);

    // Computes the average of a list of numbers.
    double average(const std::vector<double>& numbers);

   private:
    Calculator basic_calc_;  // Composition: has a Calculator instance.
};

#endif  // ADVANCED_CALCULATOR_H