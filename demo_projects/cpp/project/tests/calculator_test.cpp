#include "gtest/gtest.h"
#include "calculator.h"  // Include the header for the class we are testing.

// A test fixture can be used to share setup code between tests.
class CalculatorTest : public ::testing::Test {
   protected:
    Calculator calc;
};

TEST_F(CalculatorTest, HandlesAddition) {
    ASSERT_EQ(calc.add(2, 3), 5);
    ASSERT_EQ(calc.add(-1, 1), 0);
}

TEST_F(CalculatorTest, HandlesSubtraction) {
    ASSERT_EQ(calc.subtract(5, 2), 3);
}

TEST_F(CalculatorTest, HandlesDivision) {
    ASSERT_DOUBLE_EQ(calc.divide(10.0, 4.0), 2.5);
    // Test for exception when dividing by zero.
    ASSERT_THROW(calc.divide(1.0, 0.0), std::invalid_argument);
}

TEST_F(CalculatorTest, HandlesSignFunction) {
    // We test the positive and negative cases...
    ASSERT_EQ(calc.sign(100), 1);
    ASSERT_EQ(calc.sign(-50), -1);
    // ... but we forget to test the zero case.
}

TEST_F(CalculatorTest, HandlesTemplateSum) {
    ASSERT_EQ(calc.sum(std::vector<int>{1, 2, 3}), 6);
    ASSERT_DOUBLE_EQ(calc.sum(std::vector<double>{1.5, 2.5, 3.0}), 7.0);
    // Test the empty vector case.
    ASSERT_EQ(calc.sum(std::vector<int>{}), 0);
}