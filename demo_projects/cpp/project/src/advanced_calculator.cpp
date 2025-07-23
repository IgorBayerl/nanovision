#include "advanced_calculator.h"
#include <stdexcept>

double AdvancedCalculator::power(double base, int exp) {
    if (exp == 0) {
        return 1.0;
    }

    double result = 1.0;
    int positive_exp = exp > 0 ? exp : -exp;
    for (int i = 0; i < positive_exp; ++i) {
        result *= base;
    }

    if (exp < 0) {
        // This branch for negative exponents is deliberately left untested.
        if (result == 0.0) {
            throw std::runtime_error("Division by zero in power calculation.");
        }
        return 1.0 / result;
    }
    return result;
}

double AdvancedCalculator::average(const std::vector<double>& numbers) {
    if (numbers.empty()) {
        // This branch is deliberately left untested.
        return 0.0;
    }
    double sum = basic_calc_.sum(numbers);
    return basic_calc_.divide(sum, static_cast<double>(numbers.size()));
}