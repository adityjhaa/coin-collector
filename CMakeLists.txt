cmake_minimum_required(VERSION 3.20)
project(CoinCollector LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Add server & client subprojects
add_subdirectory(server)
add_subdirectory(client)

