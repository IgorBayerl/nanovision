
#include "math_utils.h"

namespace MathUtils {

bool is_prime(int n) {
    if (n <= 1) {
        return false;  // This branch will be deliberately missed.
    }
    for (int i = 2; i * i <= n; ++i) {
        if (n % i == 0) {
            return false;
        }
    }
    return true;
}

long long factorial(int n) {
    if (n < 0) {
        // This exception path will be missed by tests.
        throw std::invalid_argument("Factorial is not defined for negative numbers.");
    }
    if (n == 0) {
        return 1;  // This path will also be missed.
    }
    long long result = 1;
    for (int i = 1; i <= n; ++i) {
        result *= i;
    }
    return result;
}

}  // namespace MathUtils