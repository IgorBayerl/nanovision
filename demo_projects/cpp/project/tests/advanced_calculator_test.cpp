#include "gtest/gtest.h"
#include "advanced_calculator.h"

class AdvancedCalculatorTest : public ::testing::Test {
   protected:
    AdvancedCalculator adv_calc;
};

TEST_F(AdvancedCalculatorTest, Power) {
    EXPECT_DOUBLE_EQ(adv_calc.power(2.0, 3), 8.0);
    EXPECT_DOUBLE_EQ(adv_calc.power(5.0, 0), 1.0);
    // We "forget" to test negative exponents.
}

TEST_F(AdvancedCalculatorTest, Average) {
    std::vector<double> nums = {1.0, 2.0, 3.0, 4.0};
    EXPECT_DOUBLE_EQ(adv_calc.average(nums), 2.5);
    // We "forget" to test with an empty vector.
}