# project/CMakeLists.txt
cmake_minimum_required(VERSION 3.14)
project(GcovExample LANGUAGES CXX)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

enable_testing()

# ---------------------------------------------------------------------------
# Coverage flags (GCC / Clang)
# ---------------------------------------------------------------------------
set(COVERAGE_FLAGS -O0 -g --coverage)
if (CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang")
    add_compile_options(${COVERAGE_FLAGS})
    add_link_options(--coverage)
endif()

# ---------------------------------------------------------------------------
# GoogleTest
# ---------------------------------------------------------------------------
include(FetchContent)
FetchContent_Declare(
    googletest
    URL https://github.com/google/googletest/archive/refs/tags/v1.14.0.zip
)
FetchContent_MakeAvailable(googletest)

# ---------------------------------------------------------------------------
# Library under test
# ---------------------------------------------------------------------------
# Add all source files to the library.
add_library(app_lib STATIC
    src/calculator.cpp
    src/advanced_calculator.cpp
    src/utils/math_utils.cpp
)
# Add include directories for the library and its consumers.
target_include_directories(app_lib PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/utils
)

# ---------------------------------------------------------------------------
# Unit-test executable
# ---------------------------------------------------------------------------
# Add all test files to the executable.
add_executable(run_tests
    tests/calculator_test.cpp
    tests/advanced_calculator_test.cpp
    tests/math_utils_test.cpp
)

# Use the gtest_main that comes with GoogleTest – no custom main() needed.
target_link_libraries(run_tests PRIVATE app_lib GTest::gtest_main)

include(GoogleTest)
gtest_discover_tests(run_tests)