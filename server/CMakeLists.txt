add_executable(GameServer src/main.cpp)

# Winsock linking (Windows)
if (WIN32)
    target_link_libraries(GameServer PRIVATE ws2_32)
endif()
