find_package(raylib CONFIG REQUIRED)

add_executable(GameClient src/main.cpp)
target_link_libraries(GameClient PRIVATE raylib)
