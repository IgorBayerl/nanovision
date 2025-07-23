#include "calculator.h"

int Calculator::add(int a, int b) {
    return a + b;
}

int Calculator::subtract(int a, int b) {
    return a - b;
}

int Calculator::multiply(int a, int b) {
    // This function is not called by any test.
    return a * b;
}

double Calculator::divide(double a, double b) {
    if (b == 0.0) {
        throw std::invalid_argument("Division by zero is not allowed.");
    }
    return a / b;
}

int Calculator::sign(int x) {
    if (x > 0) {
        return 1;
    } else if (x < 0) {
        return -1;
    } else {
        // This branch will be deliberately missed by the tests.
        return 0;
    }
}