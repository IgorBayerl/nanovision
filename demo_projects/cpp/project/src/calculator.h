#ifndef CALCULATOR_H
#define CALCULATOR_H

#include <vector>
#include <numeric>
#include <stdexcept>

// A simple calculator class to demonstrate testing.
class Calculator {
   public:
    int add(int a, int b);
    int subtract(int a, int b);
    int multiply(int a, int b);  // This function will be intentionally untested.
    double divide(double a, double b);

    // A function with multiple branches to check coverage.
    int sign(int x);

    // A modern C++ template function.
    template <typename T>
    T sum(const std::vector<T>& numbers) {
        if (numbers.empty()) {
            return T{};  // Test the empty case.
        }
        // Use std::accumulate from <numeric>.
        return std::accumulate(numbers.begin(), numbers.end(), T{});
    }
};

#endif  // CALCULATOR_H