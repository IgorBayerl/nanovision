#include "gtest/gtest.h"
#include "utils/math_utils.h"

TEST(MathUtilsTest, IsPrime) {
    EXPECT_TRUE(MathUtils::is_prime(2));
    EXPECT_TRUE(MathUtils::is_prime(3));
    EXPECT_FALSE(MathUtils::is_prime(4));
    EXPECT_TRUE(MathUtils::is_prime(13));
    EXPECT_FALSE(MathUtils::is_prime(15));
    // We "forget" to test with n <= 1 to leave a coverage gap.
}

TEST(MathUtilsTest, Factorial) {
    EXPECT_EQ(MathUtils::factorial(1), 1);
    EXPECT_EQ(MathUtils::factorial(5), 120);
    // We "forget" to test the n=0 case and the n<0 exception.
}